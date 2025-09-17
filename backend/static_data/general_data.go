package main

import (
	wpmodels "ednevnik-backend/models/workspace"
)

// Cantons for static data
var Cantons = []wpmodels.Canton{
	{CantonCode: "ZDK", CantonName: "Zeničko-dobojski", Country: "BiH"},
	{CantonCode: "KS", CantonName: "Kanton Sarajevo", Country: "BiH"},
	{CantonCode: "TK", CantonName: "Tuzlanski", Country: "BiH"},
	{CantonCode: "USK", CantonName: "Unsko-Sanski", Country: "BiH"},
	{CantonCode: "BPK", CantonName: "Bosansko-podrinjski", Country: "BiH"},
	{CantonCode: "SBK", CantonName: "Srednjobosanski", Country: "BiH"},
	{CantonCode: "HNK", CantonName: "Hercegovačko-neretvanski", Country: "BiH"},
	{CantonCode: "ZHK", CantonName: "Zapadnohercegovački", Country: "BiH"},
	{CantonCode: "PK", CantonName: "Posavski", Country: "BiH"},
	{CantonCode: "K10", CantonName: "Kanton 10", Country: "BiH"},
}

// Classes for static data
var Classes = []wpmodels.Class{
	{ClassCode: "I"},
	{ClassCode: "II"},
	{ClassCode: "III"},
	{ClassCode: "IV"},
	{ClassCode: "V"},
	{ClassCode: "VI"},
	{ClassCode: "VII"},
	{ClassCode: "VIII"},
	{ClassCode: "IX"},
}

// Subjects for static data
var Subjects = []wpmodels.Subject{
	{SubjectCode: "BIO", SubjectName: "Biologija"},
	{SubjectCode: "BJZ", SubjectName: "Bosanski jezik i književnost"},
	{SubjectCode: "HJZ", SubjectName: "Hrvatski jezik i književnost"},
	{SubjectCode: "SJZ", SubjectName: "Srpski jezik i književnost"},
	{SubjectCode: "DMK", SubjectName: "Demokratija i ljudska prava"},
	{SubjectCode: "NJE", SubjectName: "Njemački jezik"},
	{SubjectCode: "TUR", SubjectName: "Turski jezik"},
	{SubjectCode: "FRA", SubjectName: "Francuski jezik"},
	{SubjectCode: "ARA", SubjectName: "Arapski jezik"},
	{SubjectCode: "ENG", SubjectName: "Engleski jezik"},
	{SubjectCode: "FIZ", SubjectName: "Fizika"},
	{SubjectCode: "GEO", SubjectName: "Geografija"},
	{SubjectCode: "GOB", SubjectName: "Građansko obrazovanje"},
	{SubjectCode: "HEM", SubjectName: "Hemija"},
	{SubjectCode: "HIS", SubjectName: "Historija"},
	{SubjectCode: "INF", SubjectName: "Informatika"},
	{SubjectCode: "KZ", SubjectName: "Kultura življenja"},
	{SubjectCode: "LK", SubjectName: "Likovna kultura"},
	{SubjectCode: "MK", SubjectName: "Muzička kultura"},
	{SubjectCode: "MM", SubjectName: "Matematika"},
	{SubjectCode: "TIZO", SubjectName: "Tjelesni i zdravstveni odgoj"},
	{SubjectCode: "TK", SubjectName: "Tehnička kultura"},
	{SubjectCode: "VJR", SubjectName: "Vjeronauka"},
	{SubjectCode: "MO", SubjectName: "Moja okolina"},
	{SubjectCode: "PRI", SubjectName: "Priroda"},
	{SubjectCode: "DRU", SubjectName: "Društvo"},
	{SubjectCode: "OT", SubjectName: "Osnove tehnike"},
	{SubjectCode: "SOCIO", SubjectName: "Sociologija"},
	{SubjectCode: "ELEKTRON", SubjectName: "Elektronika"},
	{SubjectCode: "ELMJER", SubjectName: "Električna mjerenja"},
	{SubjectCode: "ELKOL", SubjectName: "Električna kola"},
	{SubjectCode: "AUTO", SubjectName: "Automatika"},
	{SubjectCode: "RIP", SubjectName: "Računari i programiranje"},
	{SubjectCode: "DRM", SubjectName: "Digitalne računarske mašine"},
	{SubjectCode: "PRKS", SubjectName: "Praktična nastava"},
	{SubjectCode: "DEMOK", SubjectName: "Demokratija i ljudska prava"},
	{SubjectCode: "IDT", SubjectName: "Impulsna i digitalna tehnika"},
	{SubjectCode: "ELENERG", SubjectName: "Elektroenergetika"},
	{SubjectCode: "OSELEK", SubjectName: "Osnove elektrotehnike"},
	{SubjectCode: "TEHCRT", SubjectName: "Tehničko crtanje"},
}
