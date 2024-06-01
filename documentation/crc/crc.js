// Taken and de-obfuscated from: https://crccalc.com/

class BitOperator {
  static ['GetBool'](_0x41979d, _0x38b4b6) {
    return ((_0x41979d >> _0x38b4b6) & 1) == 1
  }
  static ['SetBool'](_0x4e2b65, _0x2be194, _0x2894cd) {
    if (_0x2894cd) return _0x4e2b65 | (1 << _0x2be194)
    else {
      let _0x1ae3ab = ~(1 << _0x2be194)
      return _0x4e2b65 & _0x1ae3ab
    }
  }
  static ['GetBit'](_0x5c8853, _0x1df56b) {
    return (_0x5c8853 >> _0x1df56b) & 1
  }
  static ['GetBits'](_0x4a68ac, _0x3409c2, _0x321235) {
    let _0x25ef16 = (1 << _0x321235) - 1
    return (_0x4a68ac >> _0x3409c2) & _0x25ef16
  }
  static ['SetBits'](_0xfc39, _0x184e47, _0x52cb68, _0x5b46ca) {
    let _0xb431ae = ~(((1 << _0x52cb68) - 1) << _0x184e47),
      _0x23bed6 =
        (_0xfc39 & _0xb431ae) | ((_0x5b46ca << _0x184e47) & ~_0xb431ae)
    return _0x23bed6
  }
}
function Stuff() {
  const _0x45cfa5 = HTMLHelper.GetTextAreaElement('textAreaInputData')
  let _0x7b3b54 = _0x45cfa5.value
  const _0xcabe19 = HTMLHelper.GetTextAreaElement('cmbbxDelimiter')
  let _0xbd1105 = _0xcabe19.value
  const _0x2ba284 = hexToByte(_0xbd1105),
    _0x3f1dd9 = hexToBytes(_0x7b3b54),
    _0x56a016 = document.getElementById('lblSourceDataHex')
  _0x56a016.innerHTML = bytesToHexP(new Uint8Array(_0x3f1dd9), 8, 8)
  const _0x5bc806 = ByteStuffing.Stuff(_0x3f1dd9, _0x2ba284),
    _0x5f020c = document.getElementById('lblProcessedDataHex')
  _0x5f020c.innerHTML = bytesToHexP(new Uint8Array(_0x5bc806), 8, 8)
}
function Unstuff() {
  const _0x3ccab1 = HTMLHelper.GetTextAreaElement('textAreaInputData')
  let _0xcf6aed = _0x3ccab1.value
  const _0x4f5bd6 = hexToBytes(_0xcf6aed),
    _0x4aeacc = document.getElementById('lblSourceDataHex')
  _0x4aeacc.innerHTML = bytesToHexP(new Uint8Array(_0x4f5bd6), 8, 8)
  const _0x1435b1 = ByteStuffing.Unstuff(_0x4f5bd6),
    _0x312e59 = document.getElementById('lblProcessedDataHex')
  _0x312e59.innerHTML = bytesToHexP(new Uint8Array(_0x1435b1), 8, 8)
}
function InitUi() {
  const _0x5eec0f = HTMLHelper.GetSelectElement('cmbbxDelimiter')
  for (let _0x45432f = 0; _0x45432f < 256; _0x45432f++) {
    let _0x1b9a6d = new Option()
    _0x1b9a6d.text = byteToHex(_0x45432f)
    _0x5eec0f.options.add(_0x1b9a6d)
  }
}
class ByteStuffing {
  static ['Stuff'](_0x1944b2, _0x10f8a3) {
    let _0x3584fc = new Array(_0x1944b2.length + 2),
      _0x275c98 = 0
    for (let _0xa24cfa = 0; _0xa24cfa < _0x1944b2.length; _0xa24cfa++) {
      const _0xb74c82 = _0x1944b2[_0xa24cfa]
      _0xb74c82 == _0x10f8a3
        ? ((_0x3584fc[_0x275c98] = _0xa24cfa + 1 - _0x275c98),
          (_0x275c98 = _0xa24cfa + 1))
        : (_0x3584fc[_0xa24cfa + 1] = _0xb74c82)
    }
    return (
      (_0x3584fc[_0x275c98] = _0x3584fc.length - _0x275c98 - 1),
      (_0x3584fc[_0x3584fc.length - 1] = _0x10f8a3),
      _0x3584fc
    )
  }
  static ['Unstuff'](_0x21b810) {
    let _0x3f689c = new Array(_0x21b810.length - 2),
      _0xc4e535 = _0x21b810[0]
    const _0x4f49e5 = _0x21b810[_0x21b810.length - 1]
    for (let _0x1e012f = 1; _0x1e012f < _0x21b810.length - 1; _0x1e012f++) {
      _0x1e012f == _0xc4e535
        ? ((_0x3f689c[_0x1e012f - 1] = _0x4f49e5),
          (_0xc4e535 += _0x21b810[_0x1e012f]))
        : (_0x3f689c[_0x1e012f - 1] = _0x21b810[_0x1e012f])
    }
    return _0x3f689c
  }
}
class Crc {
  constructor(params) {
    this.table = new Array(256)
    this.params = params
    this.mask = ((1 << this.params.hashSize) >>> 0) - 1
    this.mask = 4294967295 >>> (32 - params.hashSize)
    this.createTable()
  }
  ['calcCrc'](b) {
    let init = this.params.init
    if (this.params.refOut) {
      init = this.reverseBits(init)
    }
    if (this.params.refOut) {
      for (let i = 0; i < b.length; i++) {
        init = (this.table[(init ^ b[i]) & 255] ^ (init >>> 8)) >>> 0
        init = (init & this.mask) >>> 0
      }
    } else {
      let hashSz = this.params.hashSize - 8
      hashSz = hashSz < 0 ? 0 : hashSz
      for (let i = 0; i < b.length; i++) {
        init = this.table[((init >>> hashSz) ^ b[i]) & 255] ^ ((init << 8) >>> 0)
        init = (init & this.mask) >>> 0
      }
    }
    return (
      (init = (init ^ this.params.xorOut) >>> 0),
      (init & this.mask) >>> 0
    )
  }
  ['createTable']() {
    for (let i = 0; i < this.table.length; i++) {
      this.table[i] = this.createTableEntry(i)
    }
  }
  ['createTableEntry'](_idx) {
    let idx = _idx,
      hashSize = this.params.hashSize
    if (this.params.refIn) {
      idx = this.reverseBits(idx)
    } else {
      if (this.params.hashSize > 8) {
        idx = (idx << (hashSize - 8)) >>> 0
      }
    }
    let x31 = (1 << (hashSize - 1)) >>> 0
    for (let i = 0; i < 8; i++) {
      if ((idx & x31) != 0) {
        idx = (((idx << 1) >>> 0) ^ this.params.poly) >>> 0
      } else {
        idx = (idx << 1) >>> 0
      }
    }
    if (this.params.refIn) {
      idx = this.reverseBits(idx)
    }
    return (idx & this.mask) >>> 0
  }
  ['reverseBits'](value) {
    let hashSz = this.params.hashSize, ret = 0
    for (let i = hashSz - 1; i >= 0; i--) {
      ret = (ret | ((((value & 1) >>> 0) << i) >>> 0)) >>> 0
      value >>>= 1
    }
    return ret
  }
}
class CrcParams {
  constructor(
    _0x2a2109,
    _0x2cad11,
    _0x2441b3,
    _0x26f53b,
    _0x1f47de,
    _0x4b9a63,
    _0x46f3e8,
    _0x48512e
  ) {
    this.name = _0x2a2109
    this.hashSize = _0x2cad11
    this.poly = _0x2441b3
    this.init = _0x26f53b
    this.refIn = _0x1f47de
    this.refOut = _0x4b9a63
    this.xorOut = _0x46f3e8
    this.check = _0x48512e
  }
  static ['GetAllCrc']() {
    let _0x1dbfaf = this.Crc8.concat(this.Crc16, this.Crc32, this.Crc64)
    return _0x1dbfaf
  }
}
CrcParams.Crc8 = [
  new CrcParams('CRC-8', 8, 7, 0, false, false, 0, 244),
  new CrcParams('CRC-8/CDMA2000', 8, 155, 255, false, false, 0, 218),
  new CrcParams('CRC-8/DARC', 8, 57, 0, true, true, 0, 21),
  new CrcParams('CRC-8/DVB-S2', 8, 213, 0, false, false, 0, 188),
  new CrcParams('CRC-8/EBU', 8, 29, 255, true, true, 0, 151),
  new CrcParams('CRC-8/I-CODE', 8, 29, 253, false, false, 0, 126),
  new CrcParams('CRC-8/ITU', 8, 7, 0, false, false, 85, 161),
  new CrcParams('CRC-8/MAXIM', 8, 49, 0, true, true, 0, 161),
  new CrcParams('CRC-8/ROHC', 8, 7, 255, true, true, 0, 208),
  new CrcParams('CRC-8/WCDMA', 8, 155, 0, true, true, 0, 37),
]
CrcParams.Crc16 = [
  new CrcParams('CRC-16/ARC', 16, 32773, 0, true, true, 0, 47933),
  new CrcParams('CRC-16/AUG-CCITT', 16, 4129, 7439, false, false, 0, 58828),
  new CrcParams('CRC-16/BUYPASS', 16, 32773, 0, false, false, 0, 65256),
  new CrcParams('CRC-16/CCITT-FALSE', 16, 4129, 65535, false, false, 0, 10673),
  new CrcParams('CRC-16/CDMA2000', 16, 51303, 65535, false, false, 0, 19462),
  new CrcParams('CRC-16/DDS-110', 16, 32773, 32781, false, false, 0, 40655),
  new CrcParams('CRC-16/DECT-R', 16, 1417, 0, false, false, 1, 126),
  new CrcParams('CRC-16/DECT-X', 16, 1417, 0, false, false, 0, 127),
  new CrcParams('CRC-16/DNP', 16, 15717, 0, true, true, 65535, 60034),
  new CrcParams('CRC-16/EN-13757', 16, 15717, 0, false, false, 65535, 49847),
  new CrcParams('CRC-16/GENIBUS', 16, 4129, 65535, false, false, 65535, 54862),
  new CrcParams('CRC-16/KERMIT', 16, 4129, 0, true, true, 0, 8585),
  new CrcParams('CRC-16/MAXIM', 16, 32773, 0, true, true, 65535, 17602),
  new CrcParams('CRC-16/MCRF4XX', 16, 4129, 65535, true, true, 0, 28561),
  new CrcParams('CRC-16/MODBUS', 16, 32773, 65535, true, true, 0, 19255),
  new CrcParams('CRC-16/RIELLO', 16, 4129, 45738, true, true, 0, 25552),
  new CrcParams('CRC-16/T10-DIF', 16, 35767, 0, false, false, 0, 53467),
  new CrcParams('CRC-16/TELEDISK', 16, 41111, 0, false, false, 0, 4019),
  new CrcParams('CRC-16/TMS37157', 16, 4129, 35308, true, true, 0, 9905),
  new CrcParams('CRC-16/USB', 16, 32773, 65535, true, true, 65535, 46280),
  new CrcParams('CRC-16/X-25', 16, 4129, 65535, true, true, 65535, 36974),
  new CrcParams('CRC-16/XMODEM', 16, 4129, 0, false, false, 0, 12739),
  new CrcParams('CRC-A', 16, 4129, 50886, true, true, 0, 48901),
]
CrcParams.Crc32 = [
  new CrcParams(
    'CRC-32',
    32,
    79764919,
    4294967295,
    true,
    true,
    4294967295,
    3421780262
  ),
  new CrcParams(
    'CRC-32/BZIP2',
    32,
    79764919,
    4294967295,
    false,
    false,
    4294967295,
    4236843288
  ),
  new CrcParams(
    'CRC-32/JAMCRC',
    32,
    79764919,
    4294967295,
    true,
    true,
    0,
    873187033
  ),
  new CrcParams(
    'CRC-32/MPEG-2',
    32,
    79764919,
    4294967295,
    false,
    false,
    0,
    58124007
  ),
  new CrcParams(
    'CRC-32/POSIX',
    32,
    79764919,
    0,
    false,
    false,
    4294967295,
    1985902208
  ),
  new CrcParams(
    'CRC-32/SATA',
    32,
    79764919,
    1379029042,
    false,
    false,
    0,
    3480399848
  ),
  new CrcParams('CRC-32/XFER', 32, 175, 0, false, false, 0, 3171672888),
  new CrcParams(
    'CRC-32C',
    32,
    517762881,
    4294967295,
    true,
    true,
    4294967295,
    3808858755
  ),
  new CrcParams(
    'CRC-32D',
    32,
    2821953579,
    4294967295,
    true,
    true,
    4294967295,
    2268157302
  ),
  new CrcParams('CRC-32Q', 32, 2168537515, 0, false, false, 0, 806403967),
]
CrcParams.Crc64 = [
  new CrcParams('СRC-3/ROHC', 3, 3, 7, true, true, 0, 6),
  new CrcParams('СRC-4/ITU', 4, 3, 0, true, true, 0, 7),
  new CrcParams('СRC-5/EPC', 5, 9, 9, false, false, 0, 0),
  new CrcParams('СRC-5/ITU', 5, 21, 0, true, true, 0, 7),
  new CrcParams('СRC-5/USB', 5, 5, 31, true, true, 31, 25),
  new CrcParams('СRC-6/CDMA2000-A', 6, 39, 63, false, false, 0, 13),
  new CrcParams('СRC-6/CDMA2000-B', 6, 7, 63, false, false, 0, 59),
  new CrcParams('СRC-6/DARC', 6, 25, 0, true, true, 0, 38),
  new CrcParams('СRC-6/ITU', 6, 3, 0, true, true, 0, 6),
  new CrcParams('СRC-7', 7, 9, 0, false, false, 0, 117),
  new CrcParams('СRC-7/ROHC', 7, 79, 127, true, true, 0, 83),
  new CrcParams('СRC-10', 10, 563, 0, false, false, 0, 409),
  new CrcParams('СRC-10/CDMA2000', 10, 985, 1023, false, false, 0, 563),
  new CrcParams('СRC-11', 11, 901, 26, false, false, 0, 1443),
  new CrcParams('СRC-12/3GPP', 12, 2063, 0, false, true, 0, 3503),
  new CrcParams('СRC-12/CDMA2000', 12, 3859, 4095, false, false, 0, 3405),
  new CrcParams('СRC-12/DECT', 12, 2063, 0, false, false, 0, 3931),
  new CrcParams('СRC-13/BBC', 13, 7413, 0, false, false, 0, 1274),
  new CrcParams('СRC-14/DARC', 14, 2053, 0, true, true, 0, 2093),
  new CrcParams('СRC-15', 15, 17817, 0, false, false, 0, 1438),
  new CrcParams('СRC-15/MPT1327', 15, 26645, 0, false, false, 1, 9574),
  new CrcParams('СRC-24', 24, 8801531, 11994318, false, false, 0, 2215682),
  new CrcParams(
    'СRC-24/FLEXRAY-A',
    24,
    6122955,
    16702650,
    false,
    false,
    0,
    7961021
  ),
  new CrcParams(
    'СRC-24/FLEXRAY-B',
    24,
    6122955,
    11259375,
    false,
    false,
    0,
    2040760
  ),
  new CrcParams(
    'СRC-31/PHILIPS',
    31,
    79764919,
    2147483647,
    false,
    false,
    2147483647,
    216654956
  ),
  new CrcParams(
    'СRC-40/GSM',
    40,
    75628553,
    0,
    false,
    false,
    1099511627775,
    910907393606
  ),
]
class CrcResults {
  constructor(_0x16eaa6, _0x573223, _0x1a4113) {
    this.params = _0x16eaa6
    this.result = _0x573223
    this.lookupTable = _0x1a4113
  }
}
const _crcDict = {}
let _crcInitialized = false
function calcCrcForHexString(_0xc70893, _0x5c9196, _0x2bd74f, _0x2725c9) {
  const _0x570654 = new Uint8Array(hexToBytes(_0x2bd74f))
  calcCrc(_0xc70893, _0x570654, _0x5c9196, _0x2725c9)
}
function ProcessFile(_0x1b93b, _0x4bd118, _0x308c34) {
  calcCrc(_0x1b93b, new Uint8Array(_0x4bd118), false, _0x308c34)
}
function calcCrcForAsciiString(_0x5161ad, _0x40ef0f, _0x440215, _0x39262b) {
  const _0x46856a = new Array(0)
  for (let _0x22b6cd = 0; _0x22b6cd < _0x440215.length; _0x22b6cd++) {
    let _0x369ba8 = _0x440215.charCodeAt(_0x22b6cd)
    do {
      _0x46856a.push(_0x369ba8 % 256)
      _0x369ba8 = _0x369ba8 >>> 8
    } while (_0x369ba8 > 0)
  }
  calcCrc(_0x5161ad, new Uint8Array(_0x46856a), _0x40ef0f, _0x39262b)
}
function initCrcs() {
  if (_crcInitialized) {
    return
  }
  _crcInitialized = true
  const _0x28175e = CrcParams.GetAllCrc()
  _0x28175e.forEach((_0x538d1e) => {
    _crcDict[_0x538d1e.name.toUpperCase()] = new Crc(_0x538d1e)
  })
}
function getCrc(_0x3aef71, _0x359d62 = false) {
  const _0x4a60bc = _0x3aef71.toUpperCase().split(' ')
  initCrcs()
  const _0x325d69 = new Array()
  for (let _0x47887f in _crcDict) {
    for (let _0x5e7bdb of _0x4a60bc) {
      if (_0x359d62) {
        if (_0x47887f == _0x5e7bdb) {
          return _0x325d69.push(_crcDict[_0x47887f]), _0x325d69
        }
      } else {
        _0x47887f.indexOf(_0x5e7bdb) != -1 &&
          _0x325d69.push(_crcDict[_0x47887f])
      }
    }
  }
  return _0x325d69
}
function calcCrc(_0x3f43c2, _0x5bedbc, _0x4b342e, _0x37c8ae) {
  let _0x171a09 = getCrc(_0x3f43c2, _0x4b342e),
    _0x250752 = new Array(0)
  _0x171a09.forEach((_0xb1056) => {
    let _0x23aef6 = _0xb1056.calcCrc(_0x5bedbc)
    _0x250752.push(new CrcResults(_0xb1056.params, _0x23aef6, _0xb1056.table))
  })
  formTable(_0x250752, _0x37c8ae, _0x5bedbc)
  const _0x44d6f8 = HTMLHelper.GetInputElement('chkbxShowHex'),
    _0x15de0d = HTMLHelper.GetDivElement('sourceData')
  _0x44d6f8.checked
    ? ((_0x15de0d.style.display = 'block'), showHex(_0x5bedbc))
    : ((_0x15de0d.style.display = 'none'), showHex(new Uint8Array(0)))
}
const defaultData = [49, 50, 51, 52, 53, 54, 55, 56, 57]
function formTable(_0x2b43e3, _0x4b6186, _0x1f03d9) {
  let _0x11b9f3 = false
  if (_0x1f03d9.length == defaultData.length) {
    _0x11b9f3 = true
    for (let _0x375c4f = 0; _0x375c4f < _0x1f03d9.length; _0x375c4f++) {
      if (_0x1f03d9[_0x375c4f] != defaultData[_0x375c4f]) {
        _0x11b9f3 = false
        break
      }
    }
  }
  const _0x305fe6 = 'crcTable',
    _0xda33e = HTMLHelper.GetTableElement(_0x305fe6)
  while (_0xda33e.rows.length > 1) {
    _0xda33e.deleteRow(1)
  }
  let _0x42c372,
    _0x34be63 = '0x'
  switch (_0x4b6186) {
    case Outtypes.hex:
      ; (_0x42c372 = 16), (_0x34be63 = '0x')
      break
    case Outtypes.dec:
      ; (_0x42c372 = 10), (_0x34be63 = '')
      break
    case Outtypes.oct:
      ; (_0x42c372 = 8), (_0x34be63 = '0o')
      break
    case Outtypes.bin:
      ; (_0x42c372 = 2), (_0x34be63 = '0b')
      break
  }
  const _0x2014fa = HTMLHelper.GetTableRowElement('template')
  for (let _0x3d7f9f = 0; _0x3d7f9f < _0x2b43e3.length; _0x3d7f9f++) {
    const _0x6d76a3 = _0x2b43e3[_0x3d7f9f],
      _0x22e3db = _0x6d76a3.params.hashSize
    let _0x3de7a6
    switch (_0x4b6186) {
      case Outtypes.hex:
        _0x3de7a6 = _0x22e3db / 4
        break
      case Outtypes.dec:
        _0x3de7a6 = 0
        break
      case Outtypes.oct:
        _0x3de7a6 = Math.floor(_0x22e3db / 3) + (_0x22e3db % 3 > 0 ? 1 : 0)
        break
      case Outtypes.bin:
        _0x3de7a6 = _0x22e3db
        break
    }
    const _0x15f8e7 = _0xda33e.insertRow(-1)
    _0x15f8e7.innerHTML = _0x2014fa.innerHTML
    const _0x1d6ebf = _0x6d76a3.params.name,
      _0x3e0a8f = document.createElement('div')
    _0x2b43e3.length > 1 &&
      ((_0x3e0a8f.onclick = function() {
        getcrc(_0x1d6ebf)
      }),
        (_0x3e0a8f.style.cursor = 'pointer'),
        (_0x3e0a8f.style.color = '#003ba7'),
        (_0x3e0a8f.style.textDecoration = 'underline'))
    _0x3e0a8f.textContent = _0x1d6ebf
    _0x15f8e7.cells[0].appendChild(_0x3e0a8f)
    function _0xba7c8b(_0x9f0122) {
      return _0xeff753(_0x9f0122.toString(_0x42c372), _0x3de7a6).toUpperCase()
    }
    const _0x38212f = _0x6d76a3.result,
      _0x94a0db = _0x6d76a3.params,
      _0x53226a = _0x94a0db.check,
      _0x298929 = _0x94a0db.poly,
      _0x21c8f2 = _0x94a0db.init,
      _0x286ad5 = _0x94a0db.refIn,
      _0x160cef = _0x94a0db.refOut,
      _0x3d423b = _0x94a0db.xorOut
    _0x15f8e7.cells[1].innerHTML = _0x34be63 + _0xba7c8b(_0x38212f)
    _0x11b9f3 &&
      _0x6d76a3.result != _0x6d76a3.params.check &&
      (_0x15f8e7.cells[1].style.background = '#FFC0BC')
    _0x38212f == 0 && (_0x15f8e7.cells[1].style.background = '#ccdccc')
    _0x15f8e7.cells[2].innerHTML = _0x34be63 + _0xba7c8b(_0x53226a)
    _0x15f8e7.cells[3].innerHTML = _0x34be63 + _0xba7c8b(_0x298929)
    _0x15f8e7.cells[4].innerHTML = _0x34be63 + _0xba7c8b(_0x21c8f2)
    _0x15f8e7.cells[5].innerHTML = _0x286ad5.toString()
    _0x15f8e7.cells[6].innerHTML = _0x160cef.toString()
    _0x15f8e7.cells[7].innerHTML = _0x34be63 + _0xba7c8b(_0x3d423b)
    const _0x5f1925 = HTMLHelper.GetDivElement('crcLookupTable')
    _0x2b43e3.length == 1
      ? ((_0x5f1925.style.display = 'block'),
        showHex(_0x1f03d9),
        _0x5a2a5b(_0x2b43e3[0]))
      : ((_0x5f1925.style.display = 'none'), showHex(new Uint8Array(0)))
  }
  function _0x5a2a5b(_0x20a23c) {
    const _0x39c56d = document.getElementById('lblCrcLookupTable')
    _0x39c56d.innerHTML = ushortsToHex(
      _0x20a23c.lookupTable,
      _0x20a23c.params.hashSize
    )
  }
  function _0xeff753(_0x357a6e, _0x5dc707) {
    while (_0x357a6e.length < _0x5dc707) {
      _0x357a6e = '0' + _0x357a6e
    }
    return _0x357a6e
  }
}
let fileBytes
function dropHandler(_0x5c129e) {
  console.log('File(s) dropped')
  _0x5c129e.preventDefault()
  let _0x575946 = new FileReader()
  _0x575946.onload = function() {
    let _0xd6a486 = HTMLHelper.GetTextAreaElement('crcName').value
    console.log('Before ProcessFile')
    ProcessFile(_0xd6a486, _0x575946.result, getOutType())
    console.log('After ProcessFile')
  }
  if (_0x5c129e.dataTransfer.items) {
    let _0x1f2083 = _0x5c129e.dataTransfer.items[0].getAsFile()
    _0x575946.readAsArrayBuffer(_0x1f2083)
  } else {
    for (
      let _0x10644a = 0;
      _0x10644a < _0x5c129e.dataTransfer.files.length;
      _0x10644a++
    ) {
      console.log(
        '... file[' +
        _0x10644a +
        '].name = ' +
        _0x5c129e.dataTransfer.files[_0x10644a].name
      )
      _0x575946.readAsArrayBuffer(_0x5c129e.dataTransfer.files[_0x10644a])
    }
  }
}
function dragOverHandler(_0x51b834) {
  console.log('File(s) in drop zone')
  let _0x5490b8 = document.getElementById('drop_zone')
  _0x5490b8.style.border = '5px solid red'
  _0x51b834.preventDefault()
}
function dragLeaveHandler(_0x531114) {
  console.log('File(s) out drop zone')
  let _0x3cf635 = document.getElementById('drop_zone')
  _0x3cf635.style.border = '5px solid gray'
  _0x531114.preventDefault()
}
var Endianness
  ; (function(_0x2abb29) {
    _0x2abb29[(_0x2abb29.Undefined = 0)] = 'Undefined'
    _0x2abb29[(_0x2abb29.LittleEndian = 1)] = 'LittleEndian'
    _0x2abb29[(_0x2abb29.BigEndian = 2)] = 'BigEndian'
  })(Endianness || (Endianness = {}))
class BigEndianBitConverter {
  ['GetEndianness']() {
    return Endianness.BigEndian
  }
}
class LittleEndianBitConverter {
  ['GetEndianness']() {
    return Endianness.LittleEndian
  }
}
function hexToBytes(_0xf5c18) {
  const _0x20cb44 = new Array()
  do {
    _0xf5c18 = _0xf5c18.replace('0x', '')
  } while (_0xf5c18.search('0x') != -1)
  _0xf5c18 = _0xf5c18.replace(/[^A-Fa-f0-9]/g, '')
  console.info(_0xf5c18)
  for (let _0x37c6e7 = 0; _0x37c6e7 < _0xf5c18.length; _0x37c6e7 += 2) {
    _0x20cb44.push(parseInt(_0xf5c18.substr(_0x37c6e7, 2), 16))
  }
  return _0x20cb44
}
function hexToByte(_0x2240dd) {
  return hexToBytes(_0x2240dd)[0]
}
function bytesToHex(_0x3a76a8) {
  return bytesToHexP(_0x3a76a8, 8, 16)
}
function bytesToHexP(_0x14091b, _0x32dde3, _0x1f72f5) {
  let _0x292975 = ''
  for (let _0x2ca7c5 = 0; _0x2ca7c5 < _0x14091b.length; _0x2ca7c5++) {
    const _0x4298e8 =
      _0x14091b[_0x2ca7c5] < 0
        ? _0x14091b[_0x2ca7c5] + 256
        : _0x14091b[_0x2ca7c5]
    _0x292975 += byteToHex(_0x4298e8) + '&nbsp;'
    const _0x2d98da = _0x2ca7c5 + 1
    _0x2ca7c5 != 0 &&
      _0x2d98da % _0x32dde3 == 0 &&
      (_0x2d98da % _0x1f72f5 == 0
        ? (_0x292975 += '<br/>')
        : (_0x292975 += '&nbsp;'))
  }
  return _0x292975
}
function byteToHex(_0x2d44fb) {
  return '0x' + (_0x2d44fb >>> 4).toString(16) + (_0x2d44fb & 15).toString(16)
}
function bytesToHex2(_0x1c782c) {
  let _0x49c108 = []
  for (let _0x98916a = 0; _0x98916a < _0x1c782c.length; _0x98916a++) {
    const _0x5f34d2 =
      _0x1c782c[_0x98916a] < 0
        ? _0x1c782c[_0x98916a] + 256
        : _0x1c782c[_0x98916a]
    _0x49c108.push(
      '0x' +
      (_0x5f34d2 >>> 4).toString(16) +
      (_0x5f34d2 & 15).toString(16) +
      '&nbsp;'
    )
    const _0x2baabe = _0x98916a + 1
    _0x98916a != 0 &&
      _0x2baabe % 8 == 0 &&
      (_0x2baabe % 16 == 0 ? _0x49c108.push('<br/>') : _0x49c108.push('&nbsp;'))
  }
  let _0x24b727 = _0x49c108.join('')
  return _0x24b727
}
function showHex(_0x5bfda7) {
  const _0x58517a = document.getElementById('lblDataHex')
  _0x58517a.innerHTML = bytesToHex(_0x5bfda7)
}
function ushortsToHex(_0x380343, _0x10ea4f) {
  let _0x2daa5d = '',
    _0x2215e7 = _0x10ea4f > 8 ? 8 : 16
  for (let _0x4bf89f = 0; _0x4bf89f < _0x380343.length; _0x4bf89f++) {
    const _0x1ce5e3 = _0x380343[_0x4bf89f]
    let _0x10d367 = _0x1ce5e3.toString(16),
      _0x3db3ea = _0x10ea4f / 4 + (_0x10ea4f % 4 > 0 ? 1 : 0)
    while (_0x10d367.length < _0x3db3ea) {
      _0x10d367 = '0' + _0x10d367
    }
    _0x2daa5d += '0x' + _0x10d367 + '&nbsp;'
    const _0x4aa5ca = _0x4bf89f + 1
    _0x4bf89f != 0 &&
      _0x4aa5ca % 8 == 0 &&
      (_0x4aa5ca % 8 == 0 ? (_0x2daa5d += '<br/>') : (_0x2daa5d += '&nbsp;'))
  }
  return _0x2daa5d
}
class HTMLHelper {
  static ['GetElementById'](_0x77d3e9) {
    return document.getElementById(_0x77d3e9)
  }
  static ['GetElement'](_0x4e6328) {
    return this.GetElementById(_0x4e6328)
  }
  static ['GetAnchorElement'](_0x3ac1e2) {
    return this.GetElementById(_0x3ac1e2)
  }
  static ['GetAreaElement'](_0x306980) {
    return this.GetElementById(_0x306980)
  }
  static ['GetAudioElement'](_0x1e5b87) {
    return this.GetElementById(_0x1e5b87)
  }
  static ['GetBRElement'](_0x335bd2) {
    return this.GetElementById(_0x335bd2)
  }
  static ['GetBaseElement'](_0x3426df) {
    return this.GetElementById(_0x3426df)
  }
  static ['GetBody'](_0x16cb66) {
    return this.GetElementById(_0x16cb66)
  }
  static ['GetButtonElement'](_0xee76fc) {
    return this.GetElementById(_0xee76fc)
  }
  static ['GetCanvasElement'](_0x30a7c8) {
    return this.GetElementById(_0x30a7c8)
  }
  static ['GetDListElement'](_0x428c5a) {
    return this.GetElementById(_0x428c5a)
  }
  static ['GetDataElement'](_0x160399) {
    return this.GetElementById(_0x160399)
  }
  static ['GetDataListElement'](_0x29840e) {
    return this.GetElementById(_0x29840e)
  }
  static ['GetDetailsElement'](_0x54106f) {
    return this.GetElementById(_0x54106f)
  }
  static ['GetDialogElement'](_0x1c6325) {
    return this.GetElementById(_0x1c6325)
  }
  static ['GetDirectoryElement'](_0x31fca7) {
    return this.GetElementById(_0x31fca7)
  }
  static ['GetDivElement'](_0x288aa2) {
    return this.GetElementById(_0x288aa2)
  }
  static ['GetEmbedElement'](_0x28ccfd) {
    return this.GetElementById(_0x28ccfd)
  }
  static ['GetFieldSetElement'](_0x42cb10) {
    return this.GetElementById(_0x42cb10)
  }
  static ['GetFontElement'](_0x1eb2a6) {
    return this.GetElementById(_0x1eb2a6)
  }
  static ['GetFormElement'](_0x28a3e5) {
    return this.GetElementById(_0x28a3e5)
  }
  static ['GetFrameElement'](_0x5d1c30) {
    return this.GetElementById(_0x5d1c30)
  }
  static ['GetFrameSetElement'](_0x1a37fa) {
    return this.GetElementById(_0x1a37fa)
  }
  static ['GetHRElement'](_0x2eeee9) {
    return this.GetElementById(_0x2eeee9)
  }
  static ['GetHeadElement'](_0xc895dd) {
    return this.GetElementById(_0xc895dd)
  }
  static ['GetHeadingElement'](_0xdcea17) {
    return this.GetElementById(_0xdcea17)
  }
  static ['GetHtmlElement'](_0x2977fe) {
    return this.GetElementById(_0x2977fe)
  }
  static ['GetIFrameElement'](_0x3d620f) {
    return this.GetElementById(_0x3d620f)
  }
  static ['GetImageElement'](_0x263538) {
    return this.GetElementById(_0x263538)
  }
  static ['GetInputElement'](_0x3ce5e7) {
    return this.GetElementById(_0x3ce5e7)
  }
  static ['GetLIElement'](_0x161720) {
    return this.GetElementById(_0x161720)
  }
  static ['GetLabelElement'](_0x50eff4) {
    return this.GetElementById(_0x50eff4)
  }
  static ['GetLegendElement'](_0x3bbd03) {
    return this.GetElementById(_0x3bbd03)
  }
  static ['GetLinkElement'](_0x4faa9b) {
    return this.GetElementById(_0x4faa9b)
  }
  static ['GetMapElement'](_0x41dfe5) {
    return this.GetElementById(_0x41dfe5)
  }
  static ['GetMarqueeElement'](_0x938ff6) {
    return this.GetElementById(_0x938ff6)
  }
  static ['GetMediaElement'](_0x566f8b) {
    return this.GetElementById(_0x566f8b)
  }
  static ['GetMenuElement'](_0x3f07f0) {
    return this.GetElementById(_0x3f07f0)
  }
  static ['GetMetaElement'](_0x48f034) {
    return this.GetElementById(_0x48f034)
  }
  static ['GetMeterElement'](_0x110a2c) {
    return this.GetElementById(_0x110a2c)
  }
  static ['GetModElement'](_0x40f481) {
    return this.GetElementById(_0x40f481)
  }
  static ['GetOListElement'](_0xae4a2f) {
    return this.GetElementById(_0xae4a2f)
  }
  static ['GetObjectElement'](_0x1c2722) {
    return this.GetElementById(_0x1c2722)
  }
  static ['GetOptGroupElement'](_0x21880a) {
    return this.GetElementById(_0x21880a)
  }
  static ['GetOptionElement'](_0x118029) {
    return this.GetElementById(_0x118029)
  }
  static ['GetOrSVGElement'](_0x40cd72) {
    return this.GetElementById(_0x40cd72)
  }
  static ['GetOutputElement'](_0x2a83a3) {
    return this.GetElementById(_0x2a83a3)
  }
  static ['GetParagraphElement'](_0x598290) {
    return this.GetElementById(_0x598290)
  }
  static ['GetParamElement'](_0x2287a1) {
    return this.GetElementById(_0x2287a1)
  }
  static ['GetPictureElement'](_0x3124ff) {
    return this.GetElementById(_0x3124ff)
  }
  static ['GetPreElement'](_0x1a11a4) {
    return this.GetElementById(_0x1a11a4)
  }
  static ['GetProgressElement'](_0x99a678) {
    return this.GetElementById(_0x99a678)
  }
  static ['GetQuoteElement'](_0x2ea851) {
    return this.GetElementById(_0x2ea851)
  }
  static ['GetScriptElement'](_0x1e1d86) {
    return this.GetElementById(_0x1e1d86)
  }
  static ['GetSelectElement'](_0x1d599e) {
    return this.GetElementById(_0x1d599e)
  }
  static ['GetSlotElement'](_0x59677b) {
    return this.GetElementById(_0x59677b)
  }
  static ['GetSourceElement'](_0x53dee7) {
    return this.GetElementById(_0x53dee7)
  }
  static ['GetSpanElement'](_0xfe8d6a) {
    return this.GetElementById(_0xfe8d6a)
  }
  static ['GetStyleElement'](_0x10942d) {
    return this.GetElementById(_0x10942d)
  }
  static ['GetTableCaptionElement'](_0x591d57) {
    return this.GetElementById(_0x591d57)
  }
  static ['GetTableCellElement'](_0x5601f6) {
    return this.GetElementById(_0x5601f6)
  }
  static ['GetTableColElement'](_0x3ec215) {
    return this.GetElementById(_0x3ec215)
  }
  static ['GetTableDataCellElement'](_0xa58908) {
    return this.GetElementById(_0xa58908)
  }
  static ['GetTableElement'](_0x30731c) {
    return this.GetElementById(_0x30731c)
  }
  static ['GetTableHeaderCellElement'](_0x1c9100) {
    return this.GetElementById(_0x1c9100)
  }
  static ['GetTableRowElement'](_0x2ca950) {
    return this.GetElementById(_0x2ca950)
  }
  static ['GetTableSectionElement'](_0x1e388d) {
    return this.GetElementById(_0x1e388d)
  }
  static ['GetTemplateElement'](_0x3f654e) {
    return this.GetElementById(_0x3f654e)
  }
  static ['GetTextAreaElement'](_0x2682f7) {
    return this.GetElementById(_0x2682f7)
  }
  static ['GetTimeElement'](_0x156dc7) {
    return this.GetElementById(_0x156dc7)
  }
  static ['GetTitleElement'](_0x140baf) {
    return this.GetElementById(_0x140baf)
  }
  static ['GetTrackElement'](_0x5d37f3) {
    return this.GetElementById(_0x5d37f3)
  }
  static ['GetUListElement'](_0x12c7af) {
    return this.GetElementById(_0x12c7af)
  }
  static ['GetUnknownElement'](_0x419a47) {
    return this.GetElementById(_0x419a47)
  }
  static ['GetVideoElement'](_0x118e2e) {
    return this.GetElementById(_0x118e2e)
  }
}
var Outtypes
  ; (function(_0x5915aa) {
    _0x5915aa[(_0x5915aa.hex = 0)] = 'hex'
    _0x5915aa[(_0x5915aa.dec = 1)] = 'dec'
    _0x5915aa[(_0x5915aa.oct = 2)] = 'oct'
    _0x5915aa[(_0x5915aa.bin = 3)] = 'bin'
  })(Outtypes || (Outtypes = {}))
let reader = new FileReader(),
  customAlgoName = null,
  previousAlgoName = null
function readFile(_0xc33410) {
  if (_0xc33410 == null || _0xc33410.files.length <= 0) {
    console.error('showFile input null or files <= 0')
    return
  }
  const _0x682bfd = _0xc33410.files[0]
  reader.readAsArrayBuffer(_0x682bfd)
  reader.onload = function() {
    console.log('Reader loaded file succsesfully')
  }
  reader.onerror = function() {
    console.error('FileReader error: ' + reader.error)
  }
}
function calcFile(_0x3baefe) {
  ProcessFile(_0x3baefe, reader.result, getOutType())
}
function calcText(_0x47b758) {
  const _0x2e59a9 = HTMLHelper.GetTextAreaElement('crc'),
    _0xe4f3a = _0x2e59a9.value
  _0x2e59a9.style.backgroundColor = '#eeeeee'
  const _0x1fe97e = HTMLHelper.GetInputElement('rdbtnHex').checked,
    _0x2d4be7 = getOutType()
  console.log(_0x1fe97e)
  if (_0x1fe97e) {
  } else {
  }
}
function inputTypeOnChange() {
  const _0x628b8d = HTMLHelper.GetInputElement('rdbtnTextInputType'),
    _0x3168c0 = _0x628b8d.checked,
    _0x2d5e47 = HTMLHelper.GetParagraphElement('textInput'),
    _0x186fe5 = HTMLHelper.GetParagraphElement('fileInput')
  _0x3168c0
    ? ((_0x2d5e47.style.display = 'block'), (_0x186fe5.style.display = 'none'))
    : ((_0x2d5e47.style.display = 'none'), (_0x186fe5.style.display = 'block'))
}
function getOutType() {
  let _0x327393 = Outtypes.hex,
    _0x263268 = HTMLHelper.GetInputElement('outTyped')
  return (
    _0x263268.checked && (_0x327393 = Outtypes.dec),
    (_0x263268 = HTMLHelper.GetInputElement('outTypeo')),
    _0x263268.checked && (_0x327393 = Outtypes.oct),
    (_0x263268 = HTMLHelper.GetInputElement('outTypeb')),
    _0x263268.checked && (_0x327393 = Outtypes.bin),
    _0x327393
  )
}
function getcrc(_0x102223) {
  if (_0x102223 == null) {
    if (previousAlgoName == null) {
      alert('Please, select algo from list first')
      return
    } else {
      _0x102223 = previousAlgoName
    }
  }
  previousAlgoName = _0x102223
  const _0x74c743 =
    ' !"#$%&\'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~',
    _0x486500 = HTMLHelper.GetTextAreaElement('textAreaInputData')
  let _0x16d809 = _0x486500.value,
    _0x417dd0 = 'ascii',
    _0x479138 = HTMLHelper.GetInputElement('methodh')
  _0x479138.checked && (_0x417dd0 = 'hex')
  const _0x18119f = getOutType()
  let _0x19f3ea = _0x102223,
    _0x1b216f = true
  switch (_0x102223) {
    case 'crc8':
      ; (_0x19f3ea = 'CRC-8'), (_0x1b216f = false)
      break
    case 'crc16':
      ; (_0x19f3ea = 'CRC-16 CRC-A'), (_0x1b216f = false)
      break
    case 'crc32':
      ; (_0x19f3ea = 'CRC-32'), (_0x1b216f = false)
      break
    case 'crc64':
      ; (_0x19f3ea = 'СRC'), (_0x1b216f = false)
      break
  }
  if (_0x1b216f) {
    let _0xa82ee = HTMLHelper.GetButtonElement('btnCrcCustom')
    _0xa82ee.textContent = _0x19f3ea
    customAlgoName = _0x19f3ea
  }
  const _0x5c044d = _0x1b216f ? 'none' : 'inline-block',
    _0x227360 = _0x1b216f ? 'inline-block' : 'none'
  let _0x2c7f07 = HTMLHelper.GetButtonElement('btnCrc8')
  _0x2c7f07.style.display = _0x5c044d
  _0x2c7f07 = HTMLHelper.GetButtonElement('btnCrc16')
  _0x2c7f07.style.display = _0x5c044d
  _0x2c7f07 = HTMLHelper.GetButtonElement('btnCrc32')
  _0x2c7f07.style.display = _0x5c044d
  _0x2c7f07 = HTMLHelper.GetButtonElement('btnBackToList')
  _0x2c7f07.style.display = _0x227360
  _0x2c7f07 = HTMLHelper.GetButtonElement('btnCrcCustom')
  _0x2c7f07.style.display = _0x227360
  if (_0x417dd0 === 'hex') {
    calcCrcForHexString(_0x19f3ea, _0x1b216f, _0x16d809, _0x18119f)
  } else {
    calcCrcForAsciiString(_0x19f3ea, _0x1b216f, _0x16d809, _0x18119f)
  }
  const _0x4f1f59 = location.protocol + '//' + location.host,
    _0x280d1b =
      _0x4f1f59 +
      '/?crc=' +
      _0x16d809 +
      '&method=' +
      _0x102223 +
      '&datatype=' +
      _0x417dd0 +
      '&outtype=' +
      _0x18119f
  let _0x416935 = HTMLHelper.GetInputElement('shareurl')
  _0x416935.value = _0x280d1b
  const _0x518149 = document.getElementById('enableJsNotification')
  _0x518149.style.display = 'none'
}

