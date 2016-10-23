(function() {
    'use strict';

    angular.module('dbus-profiler', ['smart-table'])
        .controller('mainCtrl', ['$scope',function($scope) {
            $scope.ifcName = "org.freedesktop.DBus";
            $scope.switchIfc = function(name) {
                $scope.ifcName = name
            }
        }])
})()
