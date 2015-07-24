# Copyright 2014 ETH Zurich
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
"""
:mod:`ext_hdr` --- Extension header classes
===========================================
"""
# Stdlib
import binascii
import logging
import struct

# SCION
from lib.packet.packet_base import HeaderBase


class ExtensionType(object):
    """
    Constants for two types of extensions. These values are shared with L4
    protocol values, and an appropriate value is placed in next_hdr type.
    """
    HOP_BY_HOP = 0
    END_TO_END = 222


class ExtensionHeader(HeaderBase):
    """
    Abstract base class for extension headers.

    :cvar MIN_LEN:
    :type MIN_LEN: int
    :ivar next_hdr:
    :type next_hdr:
    :ivar _hdr_len:
    :type _hdr_len:
    :ivar parsed:
    :type parsed:
    """
    MIN_LEN = 8
    EXT_TYPE = None  # Type of extension (hop-by-hop or end-to-end).
    EXT_NO = None  # Number of extension.
    SUBHDR_LEN = 3 # FIXME(PSz): resize to 3
    MIN_PAYLOAD_LEN = MIN_LEN - SUBHDR_LEN

    def __init__(self, raw=None):
        """
        Initialize an instance of the class ExtensionHeader.

        :param raw:
        :type raw:
        :param _hdr_len: encoded length of extension header. The length in
                         bytes is calculated as (next_hdr + 1) * 8.
        :type _hdr_len: int
        :param next_hdr: indication of a next extension header. Must be set
                         by SCIONHeader's pack().
        :type next_hdr: int
        """
        HeaderBase.__init__(self)
        self.next_hdr = 0 
        self._hdr_len = 0
        self.payload = b"\x00" * self.MIN_PAYLOAD_LEN
        if raw is not None:
            self.parse(raw)

    def parse(self, raw):
        """
        Initialize an instance of the class ExtensionHeader.

        :param raw:
        :type raw:
        """
        assert isinstance(raw, bytes)
        dlen = len(raw)
        if dlen < self.MIN_LEN:
            logging.warning("Data too short to parse extension hdr: "
                            "data len %u", dlen)
            return
        self.next_hdr, self._hdr_len, ext_no = struct.unpack("!BBB",
            raw[:self.SUBHDR_LEN])
        assert ext_no == self.EXT_NO
        assert dlen == len(self)
        self.set_payload(raw[self.SUBHDR_LEN:])
        self.parsed = True

    def set_payload(self, payload):
        """
        Set payload and update _hdr_len.
        """
        payload_len = len(payload)
        # Length of extension must be padded to 8B.
        assert not (payload_len + self.SUBHDR_LEN) % self.MIN_LEN
        self._hdr_len = payload_len + self.SUBHDR_LEN // self.MIN_LEN
        self.payload = payload

    def pack(self):
        """

        """
        # Length of extension must be padded to 8B.
        assert not (len(self.payload) + self.SUBHDR_LEN) % self.MIN_LEN
        # next_hdr must be set for packing.
        assert self.next_hdr is not None
        return (struct.pack("!BBB", self.next_hdr, self._hdr_len, self.EXT_NO) +
                self.payload)

    def __len__(self):
        """
        Return length of extenion header in bytes.
        """
        return (self._hdr_len + 1) * self.MIN_LEN

    def __str__(self):
        """

        """
        payload_hex = binascii.hexlify(self.payload)
        return "[EH next hdr: %u, len: %u, payload: %s]" % (self.next_hdr,
                                                            len(self),
                                                            payload_hex)


class HopByHopExtension(ExtensionHeader):
    """
    Base class for hop-by-hop extensions.
    """
    EXT_TYPE = ExtensionType.HOP_BY_HOP


class EndToEndExtension(ExtensionHeader):
    """
    Base class for end-to-end extensions.
    """
    EXT_TYPE = ExtensionType.END_TO_END
