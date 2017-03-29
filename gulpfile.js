var gulp = require('gulp')
var sass = require('gulp-sass')
var uglifycss = require('gulp-uglifycss');

gulp.task('sass', function() {
  return gulp.src('static/sass/*.scss')
    .pipe(sass({
        includePaths: ['node_modules']
    }).on('error', sass.logError))
    .pipe(uglifycss())
    .pipe(gulp.dest('static/css'));
});