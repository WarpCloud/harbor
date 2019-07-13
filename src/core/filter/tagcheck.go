// Copyright 2018 Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package filter

import (
	"fmt"
	"github.com/goharbor/harbor/src/core/config"
	"net/http"
	"regexp"
	"strings"

	"github.com/astaxie/beego/context"
	"github.com/goharbor/harbor/src/common/utils/log"
)

const (
	projectTagURL  = `^/v2/((?:[a-z0-9]+(?:[._-][a-z0-9]+)*/)+)manifests/([\w][\w.-]{0,127})$`
)


// TagCheckFilter filters the deletion or creation (e.g. retag) of repo/tag requests and returns 503.
func TagCheckFilter(ctx *context.Context) {
	filterImageTag(ctx.Request, ctx.ResponseWriter)
}

func filterImageTag(req *http.Request, resp http.ResponseWriter) {

	log.Infof("begin filter Image Tag")
	if req.Method != http.MethodPut {
		return
	}

	if config.GetImageTagFilterConfig() == nil {
		return
	}

	log.Infof("begin match projectTagUrl")
	log.Infof("URL Path: %s ", req.URL.Path)
	re := regexp.MustCompile(projectTagURL)
	if ! re.MatchString(req.URL.Path) {
		return
	}

	imageFilterConfig := config.GetImageTagFilterConfig()

	for key, value := range *imageFilterConfig {

		if ! strings.HasPrefix(req.URL.Path,fmt.Sprintf("/v2/%s/", key)) {
			log.Infof("project /v2/%s/ is not in path %s, skip check...", key, req.URL.Path)
			continue
		}

		tagNameList := strings.Split(req.URL.Path, "/")
		tagName := tagNameList[len(tagNameList)-1]
		blackList := value.BlackList
		whiteList := value.WhiteList

		if matchTagBlackList(blackList, tagName) {
			log.Warningf("The image tag %s is in blackList. Any push request is prohibited.", tagName)
			resp.WriteHeader(http.StatusNotAcceptable)
			_, err := resp.Write([]byte("The image tag is in blackList. Any push request is prohibited."))
			if err != nil {
				log.Errorf("failed to write response body: %v", err)
			}
		}

		if ! matchTagWhiteList(whiteList, tagName) {
			log.Warningf("The image tag %s is not in whiteList. Any push request is prohibited.", tagName)
			resp.WriteHeader(http.StatusNotAcceptable)
			_, err := resp.Write([]byte("The image tag is not in whiteList. Any push request is prohibited."))
			if err != nil {
				log.Errorf("failed to write response body: %v", err)
			}
		}


	}

}


// matchTagBlackList checks whether a image tag is in blacklist,
// it should be blocked if  tag is in blacklist
func matchTagBlackList(blackList []string, tagName string) bool {

	if len(blackList) == 0 {
		return  false
	}

	for _, blackTag := range blackList {
		log.Infof("blackTag %s, actualTag: %s", blackTag, tagName)
		re := regexp.MustCompile(blackTag)
		if re.MatchString(tagName) {
			return true
		}
	}

	return false
}

//// matchTagWhiteList checks whether a image tag is in whiteList,
//// it should be blocked if  tag is not in blacklist
func matchTagWhiteList(whiteList []string, tagName string) bool {

	if len(whiteList) == 0 {
		return  true
	}

	for _, whiteTag := range whiteList {
		log.Infof("whiteTag %s, actualTag: %s", whiteTag, tagName)

		re := regexp.MustCompile(whiteTag)
		if re.MatchString(tagName) {
			return true
		}
	}

	return false

}
