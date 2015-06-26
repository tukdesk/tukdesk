'use strict';

angular.module("tukdesk")
    .controller("initCtrl", ["$scope", "$location", "api", function($scope, $location, api) {
        var initInfoReset = function() {
            $scope.initInfo = {
                "brandName": "",
                "email": "",
                "password": "",
                "name": ""
            }
        };

        initInfoReset();

        $scope.initInfoSubmit = function() {
            api.request.post("/base/init", $scope.initInfo).then(function() {
                initInfoReset();
                $location.path("/signin")
            }, api.responseErr());
        };
    }]);