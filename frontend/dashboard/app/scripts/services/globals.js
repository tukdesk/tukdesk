'use strict';

angular.module("tukdesk")
    .factory("globals", function($q, api) {
        var fac = {};

        // brand info
        var brandInfoReset = function() {
            fac.brand = {
                "_loaded": false
            };
        };

        brandInfoReset();

        var brandInfoSet = function(data) {
            fac.brand = data;
            fac.brand._loaded = true;
        };

        fac.brandInfoLoad = function(force) {
            var d = $q.defer();
            if (force !== true && fac.brand._loaded === true) {
                d.resolve(fac.brand);
            } else {
                api.resBrand.get()
                    .$promise.then(function(data) {
                        brandInfoSet(data);
                        d.resolve(fac.brand);
                    }, function(err) {
                        brandInfoReset();
                        d.reject(err)
                    });
            }

            return d.promise;
        };
        // brand info end

        // user info
        var userInfoReset = function() {
            fac.user = {
                "_logged": false
            }
        };

        userInfoReset();

        var userInfoSet = function(data) {
            fac.user = data;
            fac.user._logged = true;
        };

        fac.userInfoLoad = function(force) {
            var d = $q.defer();
            if (force !== true && fac.user._logged === true) {
                d.resolve(fac.user)
            } else {
                api.resProfile.get()
                    .$promise.then(function(data) {
                        userInfoSet(data);
                        d.resolve(fac.user)
                    }, function(err) {
                        userInfoReset();
                        d.reject(err)
                    });
            }

            return d.promise;
        };

        fac.userIsAgent = function() {
            return fac.user._logged && fac.user.channel.name === "_AGENT";
        };
        // user info end

        // ticket list
        fac.ticketsGlobalListReset = function() {
            fac.ticketsGlobal.list = {
                "count": 0,
                "items": []
            };
            fac.ticketsGlobal.listLoaded = false;
        };

        fac.ticketsGlobalListViewReset = function() {
            fac.ticketsGlobal.view = {
                filter: {}
            };
            fac.ticketsGlobal.viewInitialized = false;
        };

        fac.ticketsGlobalListLoad = function(forced, errCb) {
            if (forced === true || fac.ticketsGlobal.listLoaded === false) {
                return api.resTickets.list(fac.ticketsGlobal.view.filter)
                    .$promise.then(function(data) {
                        fac.ticketsGlobal.list = data;
                        fac.ticketsGlobal.listLoaded = true;
                    }, api.resourceErr(errCb))
            }
        };

        fac.ticketsGlobalListRefresh = function(errCb) {
            return fac.ticketsGlobalListLoad(true, errCb);
        };

        fac.ticketsGlobalReset = function() {
            fac.ticketsGlobal = {};
            fac.ticketsGlobalListReset();
            fac.ticketsGlobalListViewReset();
        };

        fac.ticketsGlobalReset();
        // ticket list end

        return fac;
    });
