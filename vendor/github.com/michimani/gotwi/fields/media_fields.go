package fields

type MediaField string

const (
	MediaFieldDurationMs       MediaField = "duration_ms"
	MediaFieldHeight           MediaField = "height"
	MediaFieldMediaKey         MediaField = "media_key"
	MediaFieldPreviewImageUrl  MediaField = "preview_image_url"
	MediaFieldType             MediaField = "type"
	MediaFieldUrl              MediaField = "url"
	MediaFieldWidth            MediaField = "width"
	MediaFieldPublicMetrics    MediaField = "public_metrics"
	MediaFieldNonPublicMetrics MediaField = "non_public_metrics"
	MediaFieldOrganicMetrics   MediaField = "organic_metrics"
	MediaFieldPromotedMetrics  MediaField = "promoted_metrics"
	MediaFieldAltText          MediaField = "alt_text"
	MediaFieldVariants         MediaField = "variants"
)

func (f MediaField) String() string {
	return string(f)
}

type MediaFieldList []MediaField

func (fl MediaFieldList) FieldsName() string {
	return "media.fields"
}

func (fl MediaFieldList) Values() []string {
	if fl == nil {
		return []string{}
	}

	s := []string{}
	for _, f := range fl {
		s = append(s, f.String())
	}

	return s
}
