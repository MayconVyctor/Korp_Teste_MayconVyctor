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

  successMessage: string = '';
  errorMessage: string = '';

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

      this.successMessage = '';
      this.errorMessage = '';
      
      this.stockService.createProduct(newProduct).subscribe({
        next: (savedProduct) => {
          this.successMessage = 'Product saved successfully!';
          this.loadProducts();
          this.productForm.reset({ balance: 0 });
          console.log('Product saved:', savedProduct);
          setTimeout(() => {
            this.successMessage = '';
            this.cdr.detectChanges();
          }, 3000);
        },
       error: (err) => {
          if (err.status === 500) {
            this.errorMessage = 'Error this product code already exists.';
          } else {
            this.errorMessage = 'Server error while trying to save the product.';
          }
          this.cdr.detectChanges();
        }
      });
    }
  }
}
