package loader

import "time"

const (
	PageSize              = 100
	MaxStarsCount         = 1_000_000
	MinStarsCount         = 200
	MaxConcurrentRequests = 20

	MaxRetries        = 2
	SleepTimeout      = time.Second * 90
	ErrorSleepTimeout = time.Minute * 10
)

// These are 10 next page cursors for a 100 pageSize
var Cursors = [10]string{
	"",
	"Y3Vyc29yOjEwMA==",
	"Y3Vyc29yOjIwMA==",
	"Y3Vyc29yOjMwMA==",
	"Y3Vyc29yOjQwMA==",
	"Y3Vyc29yOjUwMA==",
	"Y3Vyc29yOjYwMA==",
	"Y3Vyc29yOjcwMA==",
	"Y3Vyc29yOjgwMA==",
	"Y3Vyc29yOjkwMA==",
}

var StarsUpperBounds = []int{
	1000000,
	27104,
	17602,
	13406,
	10819,
	9083,
	7904,
	6967,
	6207,
	5616,
	5135,
	4742,
	4401,
	4106,
	3852,
	3616,
	3409,
	3235,
	3066,
	2915,
	2781,
	2663,
	2544,
	2440,
	2343,
	2247,
	2163,
	2090,
	2016,
	1952,
	1889,
	1831,
	1781,
	1728,
	1678,
	1630,
	1588,
	1547,
	1509,
	1474,
	1438,
	1403,
	1370,
	1338,
	1307,
	1277,
	1247,
	1217,
	1190,
	1164,
	1139,
	1116,
	1091,
	1068,
	1047,
	1026,
	1006,
	985,
	968,
	952,
	936,
	920,
	905,
	890,
	876,
	863,
	850,
	837,
	824,
	812,
	800,
	788,
	777,
	766,
	756,
	747,
	737,
	728,
	719,
	710,
	701,
	692,
	683,
	674,
	665,
	657,
	649,
	642,
	635,
	627,
	620,
	613,
	606,
	599,
	592,
	585,
	579,
	573,
	567,
	561,
	555,
	549,
	543,
	538,
	533,
	528,
	523,
	518,
	513,
	508,
	503,
	498,
	493,
	489,
	485,
	481,
	477,
	473,
	469,
	465,
	461,
	457,
	453,
	449,
	445,
	441,
	437,
	433,
	429,
	426,
	423,
	420,
	417,
	414,
	411,
	408,
	406,
	403,
	400,
	397,
	394,
	391,
	388,
	385,
	382,
	379,
	376,
	374,
	372,
	370,
	368,
	366,
	364,
	362,
	360,
	358,
	356,
	354,
	352,
	350,
	348,
	346,
	344,
	342,
	340,
	338,
	336,
	334,
	332,
	329,
	327,
	325,
	323,
	321,
	319,
	317,
	315,
	313,
	311,
	309,
	307,
	305,
	303,
	301,
	299,
	297,
	296,
	294,
	293,
	292,
	291,

	290,
	289,
	288,
	287,
	286,
	285,
	284,
	283,
	282,
	281,
	280,
	279,
	278,
	277,
	276,
	275,
	274,
	273,
	272,
	271,
	270,
	269,
	268,
	267,
	266,
	265,
	264,
	263,
	262,
	261,
	260,
	259,
	258,
	257,
	256,
	255,
	254,
	253,
	252,
	251,
	250,
	249,
	248,
	247,
	246,
	245,
	244,
	243,
	242,
	241,
	240,
	239,
	238,
	237,
	236,
	235,
	234,
	233,
	232,
	231,
	230,
	229,
	228,
	227,
	226,
	225,
	224,
	223,
	222,
	221,
	220,
	219,
	218,
	217,
	216,
	215,
	214,
	213,
	212,
	211,
	210,
	209,
	208,
	207,
	206,
	205,
	204,
	203,
	202,
	201,
	200,
}
