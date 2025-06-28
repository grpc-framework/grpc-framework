# Generic Go utility functions

APIs available:

* Simple logging library
    * log.go
* String Set: A list where each element appears only once
    * stringset.go
* Status Group: Like wait group, but returns error from goroutines
    * statusgroup.go
* JSON Helper Functions:
    * bodystring.go
    * jsonutils.go
* sort.StringSlice Helper functions:
    * sortedstrings.go
* UUID 16 bytes:
    * uuid.go
* SSH with SCP support:
    * ssh/
* Asynchronous HTTP Handler:
    * asynchttp/
* Super simple Unit Test assert:
    * tests/assert.go
    * tests/patch.go
* Workers producer/consumer helper function
    * workers/

# Licensing

lputils is licensed under the Apache License, Version 2.0.  See [LICENSE](https://github.com/lpabon/lputils/blob/master/LICENSE) for the full license text.
