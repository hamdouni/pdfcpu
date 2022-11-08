// Copyright 2019 The pdfcpu Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build ignore

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

var debug = flag.Bool("debug", false, "")

func main() {
	flag.Parse()

	// Generate standard.go.
	{
		w := &bytes.Buffer{}
		w.WriteString(header)
		writeWinAnsiGlyphMap(w)
		writeSymbolGlyphMap(w)
		writeZapfDingbatsGlyphMap(w)
		writeCoreFontMetrics(w)
		finish(w, "standard.go")
	}

}

func writeWinAnsiGlyphMap(w *bytes.Buffer) {
	s := `// WinAnsiGlyphMap is a glyph lookup table for CP1252 character codes.
	// See Annex D.2 Latin Character Set and Encodings.
	var WinAnsiGlyphMap = map[int]string {
	`
	writeGlyphMap(w, s, winAnsiGlyphMap)
}

func writeSymbolGlyphMap(w *bytes.Buffer) {
	s := `// SymbolGlyphMap is a glyph lookup table for Symbol character codes.
	// See Annex D.5 Symbol Set and Encoding.
	var SymbolGlyphMap = map[int]string {
	`
	writeGlyphMap(w, s, symbolGlyphMap)

}

func writeZapfDingbatsGlyphMap(w *bytes.Buffer) {
	s := `// ZapfDingbatsGlyphMap is a glyph lookup table for ZapfDingbats character codes.
	// See Annex D.6 ZapfDingbats Set and Encoding
	var ZapfDingbatsGlyphMap = map[int]string {
	`
	writeGlyphMap(w, s, zapfDingbatsGlyphMap)
}

func writeGlyphMap(w *bytes.Buffer, varDec string, m map[int]string) {
	w.WriteString(varDec)
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		fmt.Fprintf(w, "%d: \"%s\", // %#U\n", k, m[k], rune(k))
	}
	w.WriteString("}\n\n")
}

func writeCoreFontMetrics(w *bytes.Buffer) {
	s := `type fontMetrics struct {
		FBox *types.Rectangle // font box
		W    map[string]int // glyph widths
	}

	// CoreFontMetrics represents font metrics for the Adobe standard type 1 core fonts.
	var CoreFontMetrics = map[string]fontMetrics{
	`
	w.WriteString(s)
	dir := "../Core14_AFMs"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".afm") {
			continue
		}
		writeFontMetrics(w, dir, f.Name())
	}
	w.WriteString("}")
}

func writeFontBBox(w *bytes.Buffer, ss []string) {
	if len(ss) != 5 {
		panic("corrupt .afm file!")
	}
	f1, err := strconv.ParseFloat(ss[1], 64)
	if err != nil {
		log.Fatal(err)
	}
	f2, err := strconv.ParseFloat(ss[2], 64)
	if err != nil {
		log.Fatal(err)
	}
	f3, err := strconv.ParseFloat(ss[3], 64)
	if err != nil {
		log.Fatal(err)
	}
	f4, err := strconv.ParseFloat(ss[4], 64)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "types.NewRectangle(%.1f, %.1f, %.1f, %.1f),\n", f1, f2, f3, f4)
}

func writeFontMetrics(w *bytes.Buffer, dir, fileName string) {
	fmt.Fprintf(w, "\"%s\": {\n", fileName[:len(fileName)-4])
	f, err := os.Open(filepath.Join(dir, fileName))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	isHeader := true
	var headerDigested bool
	for s.Scan() {
		ss := strings.Fields(s.Text())
		if isHeader {
			switch ss[0] {
			case "FontBBox":
				writeFontBBox(w, ss)
				headerDigested = true
			case "StartCharMetrics":
				if !headerDigested {
					panic("corrupt .afm file!")
				}
				isHeader = false
				w.WriteString("map[string]int{")
			}
			continue
		}
		switch ss[0] {
		case "C":
			if len(ss) < 8 {
				panic("corrupt .afm file!")
			}
			i, err := strconv.Atoi(ss[4])
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintf(w, "\"%s\": %d, ", ss[7], i)
		case "EndCharMetrics":
			w.WriteString("},\n")
			break
		}
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
	w.WriteString("\n},\n")
}

const header = `// generated by "go run gen.go". DO NOT EDIT.

package metrics

import (	
	"github.com/hamdouni/pdfcpu/pkg/types"
)
`

func finish(w *bytes.Buffer, filename string) {
	if *debug {
		os.Stdout.Write(w.Bytes())
		return
	}
	out, err := format.Source(w.Bytes())
	if err != nil {
		log.Fatalf("format.Source: %v", err)
	}
	if err := ioutil.WriteFile(filename, out, 0660); err != nil {
		log.Fatalf("ioutil.WriteFile: %v", err)
	}
}

// See Annex D.2 Latin Character Set and Encodings
var winAnsiGlyphMap = map[int]string{
	0101: "A",
	0306: "AE",
	0301: "Aacute",
	0302: "Acircumflex",
	0304: "Adieresis",
	0300: "Agrave",
	0305: "Aring",
	0303: "Atilde",
	0102: "B",
	0103: "C",
	0307: "Ccedilla",
	0104: "D",
	0105: "E",
	0311: "Eacute",
	0312: "Ecircumflex",
	0313: "Edieresis",
	0310: "Egrave",
	0320: "Eth",
	0200: "Euro",
	0106: "F",
	0107: "G",
	0110: "H",
	0111: "I",
	0315: "Iacute",
	0316: "Icircumflex",
	0317: "Idieresis",
	0314: "Igrave",
	0112: "J",
	0113: "K",
	0114: "L",
	0115: "M",
	0116: "N",
	0321: "Ntilde",
	0117: "O",
	0214: "OE",
	0323: "Oacute",
	0324: "Ocircumflex",
	0326: "Odieresis",
	0322: "Ograve",
	0330: "Oslash",
	0325: "Otilde",
	0120: "P",
	0121: "Q",
	0122: "R",
	0123: "S",
	0212: "Scaron",
	0124: "T",
	0336: "Thorn",
	0125: "U",
	0332: "Uacute",
	0333: "Ucircumflex",
	0334: "Udieresis",
	0331: "Ugrave",
	0126: "V",
	0127: "W",
	0130: "X",
	0131: "Y",
	0335: "Yacute",
	0237: "Ydieresis",
	0132: "Z",
	0216: "Zcaron",
	0141: "a",
	0341: "aacute",
	0342: "acircumflex",
	0264: "acute",
	0344: "adieresis",
	0346: "ae",
	0340: "agrave",
	0046: "ampersand",
	0345: "aring",
	0136: "asciicircum",
	0176: "asciitilde",
	0052: "asterisk",
	0100: "at",
	0343: "atilde",
	0142: "b",
	0134: "backslash",
	0174: "bar",
	0173: "braceleft",
	0175: "braceright",
	0133: "bracketleft",
	0135: "bracketright",
	0246: "brokenbar",
	0225: "bullet",
	0143: "c",
	0347: "ccedilla",
	0270: "cedilla",
	0242: "cent",
	0210: "circumflex",
	0072: "colon",
	0054: "comma",
	0251: "copyright",
	0244: "currency",
	0144: "d",
	0206: "dagger",
	0207: "daggerdbl",
	0260: "degree",
	0250: "dieresis",
	0367: "divide",
	0044: "dollar",
	0145: "e",
	0351: "eacute",
	0352: "ecircumflex",
	0353: "edieresis",
	0350: "egrave",
	0070: "eight",
	0205: "ellipsis",
	0227: "emdash",
	0226: "endash",
	0075: "equal",
	0360: "eth",
	0041: "exclam",
	0241: "exclamdown",
	0146: "f",
	0065: "five",
	0203: "florin",
	0064: "four",
	0147: "g",
	0337: "germandbls",
	0140: "grave",
	0076: "greater",
	0253: "guillemotleft",
	0273: "guillemotright",
	0213: "guilsinglleft",
	0233: "guilsinglright",
	0150: "h",
	0055: "hyphen",
	0151: "i",
	0355: "iacute",
	0356: "icircumflex",
	0357: "idieresis",
	0354: "igrave",
	0152: "j",
	0153: "k",
	0154: "l",
	0074: "less",
	0254: "logicalnot",
	0155: "m",
	0257: "macron",
	0265: "mu",
	0327: "multiply",
	0156: "n",
	0071: "nine",
	0361: "ntilde",
	0043: "numbersign",
	0157: "o",
	0363: "oacute",
	0364: "ocircumflex",
	0366: "odieresis",
	0234: "oe",
	0362: "ograve",
	0061: "one",
	0275: "onehalf",
	0274: "onequarter",
	0271: "onesuperior",
	0252: "ordfeminine",
	0272: "ordmasculine",
	0370: "oslash",
	0365: "otilde",
	0160: "p",
	0266: "paragraph",
	0050: "parenleft",
	0051: "parenright",
	0045: "percent",
	0056: "period",
	0267: "periodcentered",
	0211: "perthousand",
	0053: "plus",
	0261: "plusminus",
	0161: "q",
	0077: "question",
	0277: "questiondown",
	0042: "quotedbl",
	0204: "quotedblbase",
	0223: "quotedblleft",
	0224: "quotedblright",
	0221: "quoteleft",
	0222: "quoteright",
	0202: "quotesinglbase",
	0047: "quotesingle",
	0162: "r",
	0256: "registered",
	0163: "s",
	0232: "scaron",
	0247: "section",
	0073: "semicolon",
	0067: "seven",
	0066: "six",
	0057: "slash",
	0040: "space",
	0243: "sterling",
	0164: "t",
	0376: "thorn",
	0063: "three",
	0276: "threequarters",
	0263: "threesuperior",
	0230: "tilde",
	0231: "trademark",
	0062: "two",
	0262: "twosuperior",
	0165: "u",
	0372: "uacute",
	0373: "ucircumflex",
	0374: "udieresis",
	0371: "ugrave",
	0137: "underscore",
	0166: "v",
	0167: "w",
	0170: "x",
	0171: "y",
	0375: "yacute",
	0377: "ydieresis",
	0245: "yen",
	0172: "z",
	0236: "zcaron",
	0060: "zero",
}

// See Annex D.5 Symbol Set and Encoding
var symbolGlyphMap = map[int]string{
	0101: "Alpha",
	0102: "Beta",
	0103: "Chi",
	0104: "Delta",
	0105: "Epsilon",
	0110: "Eta",
	0240: "Euro",
	0107: "Gamma",
	0301: "Ifraktur",
	0111: "Iota",
	0113: "Kappa",
	0114: "Lambda",
	0115: "Mu",
	0116: "Nu",
	0127: "Omega",
	0117: "Omicron",
	0106: "Phi",
	0120: "Pi",
	0131: "Psi",
	0302: "Rfraktur",
	0122: "Rho",
	0123: "Sigma",
	0124: "Tau",
	0121: "Theta",
	0125: "Upsilon",
	0241: "Upsilon1",
	0130: "Xi",
	0132: "Zeta",
	0300: "aleph",
	0141: "alpha",
	0046: "ampersand",
	0320: "angle",
	0341: "angleleft",
	0361: "angleright",
	0273: "approxequal",
	0253: "arrowboth",
	0333: "arrowdblboth",
	0337: "arrowdbldown",
	0334: "arrowdblleft",
	0336: "arrowdblright",
	0335: "arrowdblup",
	0257: "arrowdown",
	0276: "arrowhorizex",
	0254: "arrowleft",
	0256: "arrowright",
	0255: "arrowup",
	0275: "arrowvertex",
	0052: "asteriskmath",
	0174: "bar",
	0142: "beta",
	0173: "braceleft",
	0175: "braceright",
	0354: "bracelefttp",
	0355: "braceleftmid",
	0356: "braceleftbt",
	0374: "bracerighttp",
	0375: "bracerightmid",
	0376: "bracerightbt",
	0357: "braceex",
	0133: "bracketleft",
	0135: "bracketright",
	0351: "bracketlefttp",
	0352: "bracketleftex",
	0353: "bracketleftbt",
	0371: "bracketrighttp",
	0372: "bracketrightex",
	0373: "bracketrightbt",
	0267: "bullet",
	0277: "carriagereturn",
	0143: "chi",
	0304: "circlemultiply",
	0305: "circleplus",
	0247: "club",
	0072: "colon",
	0054: "comma",
	0100: "congruent",
	0343: "copyrightsans",
	0323: "copyrightserif",
	0260: "degree",
	0144: "delta",
	0250: "diamond",
	0270: "divide",
	0327: "dotmath",
	0070: "eight",
	0316: "element",
	0274: "ellipsis",
	0306: "emptyset",
	0145: "epsilon",
	0075: "equal",
	0272: "equivalence",
	0150: "eta",
	0041: "exclam",
	0044: "existential",
	0065: "five",
	0246: "florin",
	0064: "four",
	0244: "fraction",
	0147: "gamma",
	0321: "gradient",
	0076: "greater",
	0263: "greaterequal",
	0251: "heart",
	0245: "infinity",
	0362: "integral",
	0363: "integraltp",
	0364: "integralex",
	0365: "integralbt",
	0307: "intersection",
	0151: "iota",
	0153: "kappa",
	0154: "lambda",
	0074: "less",
	0243: "lessequal",
	0331: "logicaland",
	0330: "logicalnot",
	0332: "logicalor",
	0340: "lozenge",
	0055: "minus",
	0242: "minute",
	0155: "mu",
	0264: "multiply",
	0071: "nine",
	0317: "notelement",
	0271: "notequal",
	0313: "notsubset",
	0156: "nu",
	0043: "numbersign",
	0167: "omega",
	0166: "omega1",
	0157: "omicron",
	0061: "one",
	0050: "parenleft",
	0051: "parenright",
	0346: "parenlefttp",
	0347: "parenleftex",
	0350: "parenleftbt",
	0366: "parenrighttp",
	0367: "parenrightex",
	0370: "parenrightbt",
	0266: "partialdiff",
	0045: "percent",
	0056: "period",
	0136: "perpendicular",
	0146: "phi",
	0152: "phi1",
	0160: "pi",
	0053: "plus",
	0261: "plusminus",
	0325: "product",
	0314: "propersubset",
	0311: "propersuperset",
	0265: "proportional",
	0171: "psi",
	0077: "question",
	0326: "radical",
	0140: "radicalex",
	0315: "reflexsubset",
	0312: "reflexsuperset",
	0342: "registersans",
	0322: "registerserif",
	0162: "rho",
	0262: "second",
	0073: "semicolon",
	0067: "seven",
	0163: "sigma",
	0126: "sigma1",
	0176: "similar",
	0066: "six",
	0057: "slash",
	0040: "space",
	0252: "spade",
	0047: "suchthat",
	0345: "summation",
	0164: "tau",
	0134: "therefore",
	0161: "theta",
	0112: "theta1",
	0063: "three",
	0344: "trademarksans",
	0324: "trademarkserif",
	0062: "two",
	0137: "underscore",
	0310: "union",
	0042: "universal",
	0165: "upsilon",
	0303: "weierstrass",
	0170: "xi",
	0060: "zero",
	0172: "zeta",
}

// See Annex D.6 ZapfDingbats Set and Encoding
var zapfDingbatsGlyphMap = map[int]string{
	0040: "space",
	0041: "a1",
	0042: "a2",
	0043: "a202",
	0044: "a3",
	0045: "a4",
	0046: "a5",
	0047: "a119",
	0050: "a118",
	0051: "a117",
	0052: "a11",
	0053: "a12",
	0054: "a13",
	0055: "a14",
	0056: "a15",
	0057: "a16",
	0060: "a105",
	0061: "a17",
	0062: "a18",
	0063: "a19",
	0064: "a20",
	0065: "a21",
	0066: "a22",
	0067: "a23",
	0070: "a24",
	0071: "a25",
	0072: "a26",
	0073: "a27",
	0074: "a28",
	0075: "a6",
	0076: "a7",
	0077: "a8",
	0100: "a9",
	0101: "a10",
	0102: "a29",
	0103: "a30",
	0104: "a31",
	0105: "a32",
	0106: "a33",
	0107: "a34",
	0110: "a35",
	0111: "a36",
	0112: "a37",
	0113: "a38",
	0114: "a39",
	0115: "a40",
	0116: "a41",
	0117: "a42",
	0120: "a43",
	0121: "a44",
	0122: "a45",
	0123: "a46",
	0124: "a47",
	0125: "a48",
	0126: "a49",
	0127: "a50",
	0130: "a51",
	0131: "a52",
	0132: "a53",
	0133: "a54",
	0134: "a55",
	0135: "a56",
	0136: "a57",
	0137: "a58",
	0140: "a59",
	0141: "a60",
	0142: "a61",
	0143: "a62",
	0144: "a63",
	0145: "a64",
	0146: "a65",
	0147: "a66",
	0150: "a67",
	0151: "a68",
	0152: "a69",
	0153: "a70",
	0154: "a71",
	0155: "a72",
	0156: "a73",
	0157: "a74",
	0160: "a203",
	0161: "a75",
	0162: "a204",
	0163: "a76",
	0164: "a77",
	0165: "a78",
	0166: "a79",
	0167: "a81",
	0170: "a82",
	0171: "a83",
	0172: "a84",
	0173: "a97",
	0174: "a98",
	0175: "a99",
	0176: "a100",
	0241: "a101",
	0242: "a102",
	0243: "a103",
	0244: "a104",
	0245: "a106",
	0246: "a107",
	0247: "a108",
	0250: "a112",
	0251: "a111",
	0252: "a110",
	0253: "a109",
	0254: "a120",
	0255: "a121",
	0256: "a122",
	0257: "a123",
	0260: "a124",
	0261: "a125",
	0262: "a126",
	0263: "a127",
	0264: "a128",
	0265: "a129",
	0266: "a130",
	0267: "a131",
	0270: "a132",
	0271: "a133",
	0272: "a134",
	0273: "a135",
	0274: "a136",
	0275: "a137",
	0276: "a138",
	0277: "a139",
	0300: "a140",
	0301: "a141",
	0302: "a142",
	0303: "a143",
	0304: "a144",
	0305: "a145",
	0306: "a146",
	0307: "a147",
	0310: "a148",
	0311: "a149",
	0312: "a150",
	0313: "a151",
	0314: "a152",
	0315: "a153",
	0316: "a154",
	0317: "a155",
	0320: "a156",
	0321: "a157",
	0322: "a158",
	0323: "a159",
	0324: "a160",
	0325: "a161",
	0326: "a163",
	0327: "a164",
	0330: "a196",
	0331: "a165",
	0332: "a192",
	0333: "a166",
	0334: "a167",
	0335: "a168",
	0336: "a169",
	0337: "a170",
	0340: "a171",
	0341: "a172",
	0342: "a173",
	0343: "a162",
	0344: "a174",
	0345: "a175",
	0346: "a176",
	0347: "a177",
	0350: "a178",
	0351: "a179",
	0352: "a193",
	0353: "a180",
	0354: "a199",
	0355: "a181",
	0356: "a200",
	0357: "a182",
	0361: "a201",
	0362: "a183",
	0363: "a184",
	0364: "a197",
	0365: "a185",
	0366: "a194",
	0367: "a198",
	0370: "a186",
	0371: "a195",
	0372: "a187",
	0373: "a188",
	0374: "a189",
	0375: "a190",
	0376: "a191",
}
