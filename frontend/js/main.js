const leftPadding = 48;
const bottomPadding = 48;
const MaxSecond = 60;
const MaxServer = 10;
var numberServer = 0;

var iHeight = 40;
var yPosition;

var chart = d3.select('#chart svg')

function SwitchInterface(name) {
    angular.element(detail).scope().SwitchInterface(name)
}

function ShowDuration(d) {
    return (d / 1000 / 1000).toFixed(3) + "ms"
}
