import { Routes } from '@angular/router';
import {Stock} from "./components/stock/stock";

export const routes: Routes = [
    {
        path: 'stock',
        component: Stock
    },

    {
        path: '',
        redirectTo: '/stock',
        pathMatch: 'full'
    }
];
