// Copyright (c) 2019-2023, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package tools

import (
	"fmt"
	"io"
	"os"

	"github.com/sylabs/singularity/pkg/util/loop"
	"golang.org/x/sys/unix"
)

// CreateLoop associates a file to loop device and returns
// path of loop device used and a closer to close the loop device
func CreateLoop(file *os.File, offset, size uint64) (string, io.Closer, error) {
	loopDev := &loop.Device{
		MaxLoopDevices: loop.GetMaxLoopDevices(),
		Shared:         true,
		Info: &unix.LoopInfo64{
			Sizelimit: size,
			Offset:    offset,
			Flags:     unix.LO_FLAGS_AUTOCLEAR | unix.LO_FLAGS_READ_ONLY,
		},
	}
	idx := 0
	if err := loopDev.AttachFromFile(file, os.O_RDONLY, &idx); err != nil {
		return "", nil, fmt.Errorf("failed to attach image %s: %s", file.Name(), err)
	}
	return fmt.Sprintf("/dev/loop%d", idx), loopDev, nil
}
