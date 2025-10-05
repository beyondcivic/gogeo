# gogeo

```go
import "github.com/beyondcivic/gogeo/pkg/gogeo"
```

## Index

- [func Generate\(geojsonPath string, outputPath string\) \(\*geojson.FeatureCollection, error\)](<#Generate>)
- [func IsGeoJsonFile\(filePath string\) bool](<#IsGeoJsonFile>)
- [func ValidateOutputPath\(outputPath string\) error](<#ValidateOutputPath>)
- [type AppError](<#AppError>)
  - [func \(e AppError\) Error\(\) string](<#AppError.Error>)
- [type Calendar](<#Calendar>)
- [type Event](<#Event>)
- [type Geolocation](<#Geolocation>)


<a name="Generate"></a>
## func [Generate](<https://github.com:beyondcivic/gogeo/blob/main/pkg/gogeo/core.go#L10>)

```go
func Generate(geojsonPath string, outputPath string) (*geojson.FeatureCollection, error)
```

Generate generates Geo Parquet file from a geojson file with automatic type inference.

<a name="IsGeoJsonFile"></a>
## func [IsGeoJsonFile](<https://github.com:beyondcivic/gogeo/blob/main/pkg/gogeo/utils.go#L11>)

```go
func IsGeoJsonFile(filePath string) bool
```

IsCSVFile checks if a file appears to be a CSV file based on extension

<a name="ValidateOutputPath"></a>
## func [ValidateOutputPath](<https://github.com:beyondcivic/gogeo/blob/main/pkg/gogeo/utils.go#L17>)

```go
func ValidateOutputPath(outputPath string) error
```

ValidateOutputPath validates if the given path is a valid file path

<a name="AppError"></a>
## type [AppError](<https://github.com:beyondcivic/gogeo/blob/main/pkg/gogeo/error.go#L5-L10>)



```go
type AppError struct {
    // Message to show the user.
    Message string
    // Value to include with message
    Value any
}
```

<a name="AppError.Error"></a>
### func \(AppError\) [Error](<https://github.com:beyondcivic/gogeo/blob/main/pkg/gogeo/error.go#L12>)

```go
func (e AppError) Error() string
```



<a name="Calendar"></a>
## type [Calendar](<https://github.com:beyondcivic/gogeo/blob/main/pkg/gogeo/structs.go#L4-L15>)

Calendar represents a VCALENDAR component according to RFC 5545

```go
type Calendar struct {
    // Required properties
    ProdID  string `json:"prodid,omitempty"`  // Product Identifier
    Version string `json:"version,omitempty"` // iCalendar Version (should be 2.0)

    // Optional calendar properties
    CalScale string `json:"calscale,omitempty"` // Calendar scale (e.g., GREGORIAN)
    Method   string `json:"method,omitempty"`   // iTIP method (e.g., REQUEST, PUBLISH)

    // Components
    Events []Event `json:"events,omitempty"`
}
```

<a name="Event"></a>
## type [Event](<https://github.com:beyondcivic/gogeo/blob/main/pkg/gogeo/structs.go#L23-L70>)

Event represents a VEVENT component according to RFC 5545

```go
type Event struct {
    // Required properties (in most contexts)
    UID string `json:"uid,omitempty"` // Unique identifier

    // Date/Time properties
    Start    string `json:"start,omitempty"`    // DTSTART - Start date/time
    End      string `json:"end,omitempty"`      // DTEND - End date/time
    Duration string `json:"duration,omitempty"` // DURATION - Alternative to DTEND

    // Core descriptive properties
    Summary     string `json:"summary,omitempty"`     // Brief description/title
    Description string `json:"description,omitempty"` // Full description
    Location    string `json:"location,omitempty"`    // Event location

    // Optional commonly used properties
    URL        string   `json:"url,omitempty"`        // Associated URL
    Status     string   `json:"status,omitempty"`     // Event status (TENTATIVE, CONFIRMED, CANCELLED)
    Categories []string `json:"categories,omitempty"` // Event categories

    // Classification and access
    Class  string `json:"class,omitempty"`  // Access classification (PUBLIC, PRIVATE, CONFIDENTIAL)
    Transp string `json:"transp,omitempty"` // Time transparency (OPAQUE, TRANSPARENT)

    // Organizational properties
    Organizer string   `json:"organizer,omitempty"` // Event organizer
    Attendees []string `json:"attendees,omitempty"` // Event attendees

    // Scheduling properties
    Priority int `json:"priority,omitempty"` // Priority (0-9, 0=undefined)
    Sequence int `json:"sequence,omitempty"` // Revision sequence number

    // Date/Time metadata
    Created      string `json:"created,omitempty"`       // Creation date-time
    LastModified string `json:"last_modified,omitempty"` // Last modification date-time

    // Recurrence properties
    RRule        string   `json:"rrule,omitempty"`         // Recurrence rule
    RecurrenceID string   `json:"recurrence_id,omitempty"` // Recurrence identifier
    ExDates      []string `json:"exdates,omitempty"`       // Exception dates
    RDates       []string `json:"rdates,omitempty"`        // Recurrence dates

    // Other properties
    Geo       Geolocation `json:"geo,omitempty"`        // Geographic position (latitude;longitude)
    Resources []string    `json:"resources,omitempty"`  // Resources needed
    Contact   string      `json:"contact,omitempty"`    // Contact information
    RelatedTo string      `json:"related_to,omitempty"` // Related to other component
    Comment   string      `json:"comment,omitempty"`    // Comment
}
```

<a name="Geolocation"></a>
## type [Geolocation](<https://github.com:beyondcivic/gogeo/blob/main/pkg/gogeo/structs.go#L17-L20>)



```go
type Geolocation struct {
    Latitude  float64 `json:"latitude,omitempty"`
    Longitude float64 `json:"longitude,omitempty"`
}
```