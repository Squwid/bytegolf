import { HttpHeaders } from '@angular/common/http';

const httpOptions = {
    headers: new HttpHeaders({
        'Content-Type': 'application/json',
    })
};

export const Globals = {
    url: 'https://bytegolf.io',
    httpOptions,
};
