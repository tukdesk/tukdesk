'use strict';

angular.module("tukdesk")
    .controller("ticketAddCtrl", ["$scope", "api", "attachments", "broadcastEvents", function($scope, api, attachments, broadcastEvents) {
        var ticketAddInfoReset = function() {
            $scope.ticketAddInfo = {
                "channel": "_WEB",
                "email": "",
                "subject": "",
                "content": "",
                "isPublic": false,
                "attachments": []
            }
        };

        ticketAddInfoReset();

        $scope.ticketAddSubmit = function() {
            var data = angular.copy($scope.ticketAddInfo);
            data.attachmentIds = data.attachments.map(function(one) { return one.attachment.id; });
            delete data.attachments;
            
            api.resTickets.add(data)
                .$promise.then(function() {
                    ticketAddInfoReset();
                    $scope.$emit(broadcastEvents.ticketListRefresh);
                }, api.responseErr());
        };

        $scope.sizeLimit = attachments.attachmentSizeLimit;

        $scope.attachmentUpload = function(file, event, rejected) {
            if (rejected && rejected.length > 0) {
                console.log("所选文件超出规格(大小或类型)");
                return
            }

            if (!file || !file.length) {
                return
            }

            attachments.getToken(attachments.tokenUrl)
                .then(function(token) {
                    var item = {
                        progress:0
                    };

                    var progressHandler = function(event) {
                        item.progress = parseInt(100.0 * event.loaded / event.total)
                    };

                    var successHandler = function(data) {
                        item.attachment = data
                    };

                    item.uploader = attachments.uploader(file, {
                        token: token,
                        progressHandler: progressHandler,
                        successHandler: successHandler
                    });
                    $scope.ticketAddInfo.attachments.push(item);
                }, function(tokenErrResp) {
                    api.logResponseObj(tokenErrResp)
                })
        };
    }]);
