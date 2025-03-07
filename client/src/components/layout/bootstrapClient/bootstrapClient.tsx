"use client"

import { useEffect } from 'react';

export function BootstrapClient() {
    useEffect(() => {
        require('bootstrap/dist/js/bootstrap.bundle.min.js');
    }, []);

    return null;
}
