var gulp = require('gulp');
var useref = require('gulp-useref');
var uglify = require('gulp-uglify');
var gulpIf = require('gulp-if');
var cssnano = require('gulp-cssnano');
var concat = require("gulp-concat");
var shell = require("gulp-shell");

gulp.task('default', ['bindata'], function() {
})

gulp.task("bindata", shell.task(
    "go-bindata -pkg frontend js/... css/... templates/... index.html"
))
