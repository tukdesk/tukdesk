'use strict';

angular.module("tukdesk")
    .controller("ticketAddCtrl", ["$scope", "api", "broadcastEvents", function($scope, api, broadcastEvents) {
        var ticketAddInfoReset = function() {
            $scope.ticketAddInfo = {
                "channel": "_WEB",
                "email": "",
                "subject": "",
                "content": "",
                "isPublic": false
            }
        };

        ticketAddInfoReset();

        $scope.ticketAddSubmit = function() {
            api.resTickets.add($scope.ticketAddInfo)
                .$promise.then(function() {
                    ticketAddInfoReset();
                    $scope.$emit(broadcastEvents.ticketListRefresh);
                }, api.resourceErr());
        };
    }]);
