'use strict';

angular.module("tukdesk")
    .factory("attachments", function($q, $http, api, utils, Upload) {
        var fac = {};

        var token = {
            token: "",
            expired: 0
        };

        var tokenExpired = function() {
            return token.expired < utils.timestampNow();
        };

        fac.tokenUrl = api.apiUrl("/attachments/token");

        fac.getToken = function(url) {
            var d = $q.defer();
            if (tokenExpired() === false) {
                d.resolve(token.token);
            } else {
                $http.get(url).then(function(resp) {
                    var data = resp.data;
                    token.token = data.token;
                    token.expired = utils.timestampNow() + data.expiration - 300;
                    d.resolve(token.token)
                }, function(errResp) {
                    d.reject(errResp);
                });
            }

            return d.promise;
        };

        fac.uploadConfig = {
            uploadUrl: api.apiUrl("/attachments/upload"),
            progressHandler: function(event) {
//                var progressPercentage = parseInt(100.0 * event.loaded / event.total);
//                console.log('progress: ' + progressPercentage + '% ');
                console.log(event.config);
            },
            successHandler: function(data) { console.log(data); },
            errorHandler: api.logHTTPResp
        };

        fac.attachmentSizeLimit = 2097152;

        fac.uploader = function(file, conf) {
            var config = angular.copy(fac.uploadConfig);
            if (conf) {
                angular.extend(config, conf)
            }

            return Upload.upload({
                url: config.uploadUrl,
                file: file,
                fileFormDataName: "file",
                fields: {
                    token: config.token
                }
            }).progress(config.progressHandler)
                .success(config.successHandler)
                .error(config.errorHandler);
        };

        return fac;
    });
