package main

import (
	wpmodels "ednevnik-backend/models/workspace"
)

// Semesters for static data
var Semesters = []wpmodels.Semester{
	{SemesterCode: "1POL", SemesterName: "Prvo polugodište", ProgressLevel: 1},
	{SemesterCode: "2POL", SemesterName: "Drugo polugodište", ProgressLevel: 2},
}

// NPPSemester for static data
var NPPSemester = []wpmodels.NPPSemester{
	{NPPCode: "BOS", SemesterCode: "1POL", StartDate: "2023-09-01", EndDate: "2024-12-31"},
	{NPPCode: "BOS", SemesterCode: "2POL", StartDate: "2024-02-01", EndDate: "2024-06-01"},
	{NPPCode: "HRV", SemesterCode: "1POL", StartDate: "2023-09-01", EndDate: "2024-12-20"},
	{NPPCode: "HRV", SemesterCode: "2POL", StartDate: "2024-01-15", EndDate: "2024-06-01"},
}
