'use strict';

module.exports = function (grunt) {
    // Load grunt tasks automatically
    require('load-grunt-tasks')(grunt);

    // Time how long tasks take. Can help when optimizing build times
    require('time-grunt')(grunt);

    grunt.initConfig({
        tudesk: {
            app: require('./bower.json').appPath || 'app',
            dist: 'dist'
        },
        clean: {
            dist: {
                files: [{
                    dot: true,
                    src: [
                        '.tmp',
                        '<%= tudesk.dist %>/*',
                        '!<%= tudesk.dist %>/.git*'
                    ]
                }]
            },
            server: '.tmp'
        },
        less: {
            dist: {
                options: {
                    cleancss: true,
                    report: 'min'
                },
                files: [
                    {
                        expand: true,
                        cwd: "<%= tudesk.app %>/assets/less/tukdesk",
                        src: "tukdesk.less",
                        dest: "<%= tudesk.dist %>/styles",
                        ext: ".css"
                    }
                ]
            },
            dev: {
                options: {
                    cleancss: true,
                    report: 'min'
                },
                files: [
                    {
                        expand: true,
                        cwd: "<%= tudesk.app %>/assets/less/tukdesk",
                        src: "tukdesk.less",
                        dest: "<%= tudesk.app %>/styles",
                        ext: ".css"
                    }
                ]
            }
        },

        // Reads HTML for usemin blocks to enable smart builds that automatically
        // concat, minify and revision files. Creates configurations in memory so
        // additional tasks can operate on them
        useminPrepare: {
            html: '<%= tudesk.app %>/index.html',
            options: {
                dest: '<%= tudesk.dist %>',
                flow: {
                    html: {
                        steps: {
                            js: ['concat', 'uglifyjs'],
                            css: ['cssmin']
                        },
                        post: {}
                    }
                }
            }
        },

        // Performs rewrites based on rev and the useminPrepare configuration
        usemin: {
            html: ['<%= tudesk.dist %>/{,*/}*.html'],
            css: ['<%= tudesk.dist %>/assets/css/{,*/}*.css'],
            options: {
                assetsDirs: ['<%= tudesk.dist %>']
            }
        },

        // The following *-min tasks produce minified files in the dist folder
        cssmin: {
            options: {
                relativeTo: '<%= tudesk.app %>',
                processImport: true,
                noAdvanced: true
            }
        },

        // Renames files for browser caching purposes
        rev: {
            dist: {
                files: {
                    src: [
                        '<%= tudesk.dist %>/scripts/{,*/}*.js',
                        '<%= tudesk.dist %>/assets/css/{,*/}*.css',
                        '<%= tudesk.dist %>/styles/fonts/*'
                    ]
                }
            }
        },

        // Copies remaining files to places other tasks can use
        copy: {
            dist: {
                files: [{
                    expand: true,
                    dot: true,
                    cwd: '<%= tudesk.app %>',
                    dest: '<%= tudesk.dist %>',
                    src: [
                        'index.html',
                        'views/{,*/}*.html',
                        'fonts/**',
                        'assets/**'
                    ]
                }]
            },
            styles: {
                expand: true,
                cwd: '<%= tudesk.app %>/assets/css',
                dest: '.tmp/assets/css',
                src: '{,*/}*.css'
            },
            dev: {

            }
        },

        // Run some tasks in parallel to speed up the build process
        concurrent: {
            server: [
                'copy:styles'
            ],
            test: [
                'copy:styles'
            ],
            dist: [
                'copy:styles',
                'copy:dist',
            ]
        },
        uglify: {
            options: {
                mangle: false
            }
        }

    });

    grunt.registerTask('build', [
        'clean:dist',
        'less:dist',
        'useminPrepare',
        'concurrent:dist',
        'concat',
        'copy:dist',
        'cssmin',
        'uglify',
        'rev',
        'usemin'
    ]);

    grunt.registerTask('dev', [
        'clean:dist',
        'less:dev',
    ]);

    grunt.registerTask('default', [
        'build'
    ]);
};
