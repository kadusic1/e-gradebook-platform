package main

import (
	wpmodels "ednevnik-backend/models/workspace"
)

// Courses for static data
var Courses = []wpmodels.Course{
	{CourseCode: "AT", CourseName: "Arhitektonski tehničar", CourseDuration: "IV"},
	{CourseCode: "ELEL", CourseName: "Elektrotehničar elektronike", CourseDuration: "IV"},
	{CourseCode: "ENERG", CourseName: "Elektrotehničar energetike", CourseDuration: "IV"},
	{CourseCode: "GEO", CourseName: "Geodetski tehničar (geometar)", CourseDuration: "IV"},
	{CourseCode: "MASTK", CourseName: "Mašinski tehničar za kompjutersko upravljanje mašinama", CourseDuration: "IV"},
	{CourseCode: "METT", CourseName: "Metalurški tehničar", CourseDuration: "IV"},
	{CourseCode: "RTIA", CourseName: "Elektrotehničar računarske tehnike i automatike", CourseDuration: "IV"},
	{CourseCode: "ST", CourseName: "Tehničar drumskog saobraćaja", CourseDuration: "IV"},
}
