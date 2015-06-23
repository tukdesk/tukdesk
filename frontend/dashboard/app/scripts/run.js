'use strict';

angular.module("tukdesk").run(function($rootScope, api, globals, utils) {
    console.log("running");

    api.setDefaultHeader();
    $rootScope.globals = globals;
    $rootScope.utils = utils;
});
