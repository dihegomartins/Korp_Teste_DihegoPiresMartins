import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterOutlet } from '@angular/router';
import { NotaListaComponent } from './components/nota-lista/nota-lista';
// 1. Importe o novo componente aqui:
import { ProdutoCadastroComponent } from './components/produto-cadastro/produto-cadastro';

@Component({
  selector: 'app-root',
  standalone: true,
  // 2. Adicione ele aqui no array de imports:
  imports: [
    CommonModule,
    RouterOutlet,
    NotaListaComponent,
    ProdutoCadastroComponent, // <-- ESSENCIAL!
  ],
  templateUrl: './app.html',
  styleUrl: './app.scss',
})
export class AppComponent {
  title = 'faturamento-app';
}
