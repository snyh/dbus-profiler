var gulp = require('gulp');
var useref = require('gulp-useref');
var uglify = require('gulp-uglify');
var gulpIf = require('gulp-if');
var cssnano = require('gulp-cssnano');
var concat = require("gulp-concat");

gulp.task('default', ['dist'], function() {
})


gulp.task('js', ['jslib'], function() {
    return gulp.src('js/*.js')
        .pipe(concat("main.min.js"))
//        .pipe(uglify())
        .pipe(gulp.dest('dist/js'));
});
gulp.task('jslib', function() {
    return gulp.src('js/lib/*.js')
        .pipe(gulp.dest('dist/js/lib'));
});
gulp.task('csslib', function() {
    return gulp.src('css/lib/*.css')
        .pipe(gulp.dest('dist/css/lib'));
});


gulp.task('dist', ['csslib', 'js'], function() {
    //    return gulp.src(['*.html', "js/*.js"])
    return gulp.src("app/*.html")
        .pipe(useref())
        .pipe(gulp.dest('dist'))
});
