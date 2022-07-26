# Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
# Use of this source code is governed by a license that can be
# found in the LICENSE file.

#!/bin/sh
grep -vf scripts/resources/exclude-from-code-coverage.txt coverage.out > coverage.out.tmp && mv coverage.out.tmp coverage.out
