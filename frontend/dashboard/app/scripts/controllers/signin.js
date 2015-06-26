'use strict';

angular.module("tukdesk")
    .controller("signinCtrl", ["$scope", "$location", "api", function($scope, $location, api) {
        console.log("signin ctrl");

        var signinInfoReset = function() {
            $scope.signinInfo = {
                "password": ""
            }
        };

        signinInfoReset();

        $scope.signinInfoSubmit = function() {
            api.request.post("/base/signin", $scope.signinInfo)
                .success(function (data) {
                    api.setAuthHeader(data.token);
                    $location.path("/");
                })
                .error(function (data, status, headers) {
                    api.logHTTPResp(data, status, headers)
                })
        };
    }]);
