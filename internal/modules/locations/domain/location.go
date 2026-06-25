package domain

type Components struct {
	CountryCode    string
	CountryName    string
	AdminAreaCode  *string
	AdminAreaName  *string
	LocalityCode   *string
	LocalityName   *string
	PostalCode     *string
	SettlementName *string
	SettlementType *string
	StreetName     *string
}

const (
	AdminAreaTypeState        = "state"
	AdminAreaTypeMunicipality = "municipality"
)

type Country struct {
	CountryCode   string
	Name          string
	Label         string
	Source        string
	SourceVersion string
	Score         float64
}

type AdminArea struct {
	ID            string
	CountryCode   string
	Code          string
	Name          string
	Type          string
	ParentCode    *string
	Label         string
	Source        string
	SourceVersion string
	Score         float64
}

type Locality struct {
	ID            string
	CountryCode   string
	AdminAreaCode string
	Code          string
	Name          string
	Type          string
	Label         string
	Source        string
	SourceVersion string
	Score         float64
}

type Settlement struct {
	ID            string
	CountryCode   string
	AdminAreaCode string
	LocalityCode  *string
	PostalCode    *string
	Name          string
	Type          string
	Label         string
	Source        string
	SourceVersion string
	Score         float64
}

type PostalCodeMatch struct {
	PostalCode    string
	Label         string
	Components    Components
	Source        string
	SourceVersion string
	Score         float64
}

type AddressSuggestion struct {
	Label         string
	Components    Components
	Source        string
	SourceVersion string
	Score         float64
}

type CountryFilter struct {
	Query string
	Limit int32
}

type AdminAreaFilter struct {
	CountryCode string
	ParentCode  string
	Type        string
	Query       string
	Limit       int32
}

type LocalityFilter struct {
	CountryCode   string
	AdminAreaCode string
	Query         string
	Limit         int32
}

type SettlementFilter struct {
	CountryCode   string
	AdminAreaCode string
	LocalityCode  string
	PostalCode    string
	Query         string
	Limit         int32
}

type PostalCodeFilter struct {
	CountryCode      string
	PostalCodePrefix string
	Limit            int32
}

type AddressSuggestionFilter struct {
	CountryCode string
	Query       string
	PostalCode  string
	Limit       int32
}
