import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Product {
  code: string;
  description: string;
  balance: number;
}

export interface InvoiceItem {
  id?: number;
  invoice_id?: number;
  product_code: string;
  product_description?: string;
  quantity: number;
}

export interface Invoice {
  id?: number;
  status: string;
  items: InvoiceItem[];
}

export interface AiAnalysisResponse {
  invoice_id: number;
  status: string;
  ai_analysis: string;
}

@Injectable({
  providedIn: 'root'
})
export class InvoiceService {
  private billingUrl = 'http://localhost:8082/invoices';
  private inventoryUrl = 'http://localhost:8081/products'; 

  constructor(private http: HttpClient) {}

  getProducts(): Observable<Product[]> {
    return this.http.get<Product[]>(this.inventoryUrl);
  }

  getInvoices(): Observable<any> {
    return this.http.get<any>(this.billingUrl);
  }

  createInvoice(invoice: Invoice): Observable<any> {
    return this.http.post<any>(this.billingUrl, invoice);
  }

  analyzeInvoice(id: number): Observable<AiAnalysisResponse> {
    return this.http.get<AiAnalysisResponse>(`${this.billingUrl}/${id}/analysis`);
  }

  printInvoice(id: number): Observable<void> {
    return this.http.put<void>(`${this.billingUrl}/${id}/print`, {});
  }
}