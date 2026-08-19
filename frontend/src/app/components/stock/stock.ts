import { Component, OnInit, ChangeDetectorRef} from '@angular/core';
import { StockService, Product } from '../../services/stockService';

@Component({
  selector: 'app-stock',
  imports: [],
  templateUrl: './stock.html',
  styleUrl: './stock.css',
})
export class Stock implements OnInit {

  
  products: Product[] = [];

  constructor(private stockService: StockService, private cdr: ChangeDetectorRef) {}

  ngOnInit(): void {
    this.loadProducts();
  }

  loadProducts(): void {
    this.stockService.getProducts().subscribe({
      next: (data) => {
        this.products = data;
        this.cdr.detectChanges();
        console.log('Produtos carregados do Go:', this.products);
      },
      error: (err) => {
        console.error('Deu erro ao buscar no back-end:', err);
      }
    });
  }
}
