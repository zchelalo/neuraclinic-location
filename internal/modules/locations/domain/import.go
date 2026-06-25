package domain

type DataSource struct {
	Key         string
	Name        string
	Version     string
	License     string
	Attribution string
	URL         string
	CountryCode string
	CountryName string
}

type SEPOMEXRow struct {
	PostalCode       string
	SettlementName   string
	SettlementType   string
	Municipality     string
	State            string
	City             string
	StateCode        string
	MunicipalityCode string
	SettlementCode   string
	Zone             string
	CityCode         string
}

type INEGIEntityRow struct {
	GeoCode           string
	StateCode         string
	StateName         string
	StateAbbreviation string
}

type INEGIMunicipalityRow struct {
	GeoCode           string
	StateCode         string
	StateName         string
	StateAbbreviation string
	MunicipalityCode  string
	MunicipalityName  string
	HeadCode          string
	HeadName          string
}

type ImportSummary struct {
	RowsRead     int
	RowsImported int
	DataSources  int
	Countries    int
	AdminAreas   int
	Localities   int
	PostalCodes  int
	Settlements  int
}
