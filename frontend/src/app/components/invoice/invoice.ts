import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
import { FormGroup, FormControl, Validators, ReactiveFormsModule } from '@angular/forms';
import { InvoiceService, Invoice as InvoiceModel, AiAnalysisResponse, Product } from '../../services/invoiceService';

@Component({
  selector: 'app-invoice',
  imports: [ReactiveFormsModule],
  templateUrl: './invoice.html',
  styleUrl: './invoice.css',
})
export class Invoice implements OnInit {
  invoices: InvoiceModel[] = [];
  products: Product[] = [];
  successMessage: string = '';
  errorMessage: string = '';

  invoiceForm = new FormGroup({
    ProductCode: new FormControl('', Validators.required),
    Quantity: new FormControl(1, [Validators.required, Validators.min(1)])
  });

  aiAnalysisResult: string | null = null;
  isAnalyzing: boolean = false;
  analyzingId: number | null = null;
  isPrintingId: number | null = null;

  constructor(
    private invoiceService: InvoiceService, 
    private cdRef: ChangeDetectorRef
  ) {}

  ngOnInit(): void {
    this.loadInvoices();
    this.loadProducts()
  }

  loadInvoices(): void {
    this.invoiceService.getInvoices().subscribe({
      next: (data: any) => {
        this.invoices = data.invoices;
        this.cdRef.detectChanges();
      },
      error: (err) => console.error('Error fetching invoices:', err)
    });
  }

  loadProducts(): void {
    this.invoiceService.getProducts().subscribe({
      next: (data: any) => {
        this.products = data.products || data; 
        this.cdRef.detectChanges();
      },
      error: (err) => console.error('Error fetching products:', err)
    });
  }

  onSubmit(): void {
    if (this.invoiceForm.valid) {
      const newInvoice: InvoiceModel = {
        id: 0,
        status: 'OPEN', 
        items: [
          {
            id: 0,
            invoice_id: 0,
            product_code: this.invoiceForm.value.ProductCode ?? '',
            quantity: this.invoiceForm.value.Quantity ?? 1
          }
        ]
      };

      this.invoiceService.createInvoice(newInvoice).subscribe({
        next: () => {
          this.successMessage = 'Invoice issued successfully!';
          this.invoiceForm.reset({ Quantity: 1 });
          this.loadInvoices();
          setTimeout(() => { this.successMessage = ''; this.cdRef.detectChanges(); }, 6000);
        },
        error: (err) => {
          this.errorMessage = 'Error issuing invoice.';
          console.error(err);
          this.cdRef.detectChanges();
        }
      });
    }
  }

    printInvoice(id: number | undefined): void {
      if (!id) return;
      
      this.isPrintingId = id;

      this.invoiceService.printInvoice(id).subscribe({
        next: () => {
          this.loadInvoices();
          this.isPrintingId = null;
          this.cdRef.detectChanges();
        },
        error: (err) => {
          console.error('Error printing invoice:', err);
          this.isPrintingId = null;
          this.cdRef.detectChanges();
        }
      });

    }

    analyzeWithAI(id: number | undefined): void {
      if (!id) return;

      this.isAnalyzing = true;
      this.analyzingId = id;
      this.aiAnalysisResult = null;

      this.invoiceService.analyzeInvoice(id).subscribe({
        next: (res: AiAnalysisResponse) => {
          this.aiAnalysisResult = res.ai_analysis;
          this.isAnalyzing = false;
          this.analyzingId = null;
          this.cdRef.detectChanges();
        },
        error: (err) => {
          console.error('Error analyzing invoice:', err);
          this.aiAnalysisResult = 'Error analyzing invoice.';
          this.isAnalyzing = false;
          this.analyzingId = null;
          this.cdRef.detectChanges();
        }
      });
  }
}