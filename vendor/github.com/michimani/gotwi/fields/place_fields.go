package fields

type PlaceField string

const (
	PlaceFieldContainedWithin PlaceField = "contained_within"
	PlaceFieldCountry         PlaceField = "country"
	PlaceFieldCountryCode     PlaceField = "country_code"
	PlaceFieldFullName        PlaceField = "full_name"
	PlaceFieldGeo             PlaceField = "geo"
	PlaceFieldID              PlaceField = "id"
	PlaceFieldName            PlaceField = "name"
	PlaceFieldPlaceType       PlaceField = "place_type"
)

func (f PlaceField) String() string {
	return string(f)
}

type PlaceFieldList []PlaceField

func (fl PlaceFieldList) FieldsName() string {
	return "place.fields"
}

func (fl PlaceFieldList) Values() []string {
	if fl == nil {
		return []string{}
	}

	s := []string{}
	for _, f := range fl {
		s = append(s, f.String())
	}

	return s
}
