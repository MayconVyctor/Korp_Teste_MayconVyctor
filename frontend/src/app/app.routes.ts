import { Routes } from '@angular/router';
import {Stock} from "./components/stock/stock";
import {Invoice} from "./components/invoice/invoice";

export const routes: Routes = [
    {
        path: 'stock',
        component: Stock
    },

    {
        path: '',
        redirectTo: '/stock',
        pathMatch: 'full'
    },

    { 
        path: 'invoice', 
        component: Invoice 
    },
    
];
