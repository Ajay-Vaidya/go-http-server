package homepage

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"reflect"
	"slices"
	"strings"

	"go_web_server/util"
)

type HomepageData struct {
	PageTitle  string
	ApiRoutes  []string
	HtmlRoutes []string
}

func convertRouteMapKeysToStringSlice(values []reflect.Value) []string {
	var stringsSlice []string
	for _, v := range values {
		// Check if the underlying type of v is a string before conversion
		if v.Kind() == reflect.String {
			switch v.String() {
			case "/favicon.ico":
				// Do nothing for these values
			default:
				// Append value to slice for all other string values
				stringsSlice = append(stringsSlice, v.String())
			}
		}
	}
	slices.Sort(stringsSlice)
	return stringsSlice
}

func getAvailableRoutes(w http.ResponseWriter, mux *http.ServeMux) ([]string, []string) {
	v := reflect.ValueOf(mux).Elem()
	routeMap := v.FieldByName("m")
	routes := convertRouteMapKeysToStringSlice(routeMap.MapKeys())
	apiRoutes := util.Filter(routes, func(route string) bool {
		return strings.Contains(route, "/api")
	})
	htmlRoutes := util.Filter(routes, func(route string) bool {
		return !strings.Contains(route, "/api")
	})
	return apiRoutes, htmlRoutes
}

func renderHTMLTemplate(data HomepageData, w http.ResponseWriter) error {
	tmpl, err := template.ParseFiles("homepage/homepage.tmpl")
	if err != nil {
		log.Print(err)
		fmt.Fprintf(w, "Failed to Parse Template!\n")
		return err
	}
	tmplErr := tmpl.Execute(w, data)
	if tmplErr != nil {
		log.Print(tmplErr)
		fmt.Fprintf(w, "Failed to Execute Template!\n")
		return tmplErr
	}
	return nil
}

func HomePage(contentType string, mux *http.ServeMux) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Endpoint Hit: homePage")
		apiRoutes, htmlRoutes := getAvailableRoutes(w, mux)
		data := HomepageData{PageTitle: "Homepage", ApiRoutes: apiRoutes, HtmlRoutes: htmlRoutes}
		util.ResponseBody(data, contentType, w, renderHTMLTemplate)
	}
}
