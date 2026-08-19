import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';  

export interface Product {
  code: string;
  description: string;
  balance: number;
}

@Injectable({
  providedIn: 'root',
})
export class StockService {

  private apiUrl = 'http://localhost:8081/products';
  constructor(private http: HttpClient) {}

  getProducts(): Observable<Product[]> {
    return this.http.get<Product[]>(this.apiUrl);
  }
}
