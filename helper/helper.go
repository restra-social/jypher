package helper


func IDfilter(ftype string, id string) string {

	switch ftype {
	case "urn":
		// filter urn:uuid:95575f5f-feb1-459b-819f-07ac602e6f6b into 95575f5f-feb1-459b-819f-07ac602e6f6b
		return id[9:len(id)]
	}
	return id
}