(function() {
    'use strict';
    angular.module('dbus-profiler')
        .filter('sum', function() {
            return function(d) {
                var sum = 0
                angular.forEach(d, function(i) { sum = sum + i})
                return sum
            }
        })
        .filter('percentage', ['$filter', function($filter) {
            return function(input, all, decimals) {
                return $filter('number')(input * 100.0 / all, decimals) + '%';
            }
        }]);
})()
