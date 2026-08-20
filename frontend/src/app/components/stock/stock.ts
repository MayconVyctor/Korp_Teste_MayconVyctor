import { Component, OnInit, ChangeDetectorRef} from '@angular/core';
import { StockService, Product } from '../../services/stockService';
import { FormGroup, FormControl, Validators, ReactiveFormsModule } from '@angular/forms';

@Component({
  selector: 'app-stock',
  imports: [ReactiveFormsModule],
  templateUrl: './stock.html',
  styleUrl: './stock.css',
})
export class Stock implements OnInit {

  
  products: Product[] = [];

  productForm = new FormGroup({
    code: new FormControl('', Validators.required),
    description: new FormControl('', Validators.required),
    balance: new FormControl(0, [Validators.required, Validators.min(0)])
  });

  constructor(
    private stockService: StockService, 
    private cdr: ChangeDetectorRef) 
    {}

  ngOnInit(): void {
    this.loadProducts();
  }

  loadProducts(): void {
    this.stockService.getProducts().subscribe({
      next: (data) => {
        this.products = data;
        this.cdr.detectChanges();
        console.log('Products loaded:', this.products);
      },
      error: (err) => {
        console.error('Error loading products:', err);
      }
    });
  }
  onSubmit(): void {
    if (this.productForm.valid) {
      const newProduct = this.productForm.value as Product;
      
      this.stockService.createProduct(newProduct).subscribe({
        next: (savedProduct) => {
          console.log('Product saved:', savedProduct);
          this.loadProducts();
          this.productForm.reset({ balance: 0 }); 
        },
        error: (err) => {
          console.error('Error saving the product:', err);
        }
      });
    }
  }
}
