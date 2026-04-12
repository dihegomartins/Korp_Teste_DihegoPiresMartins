export interface NotaFiscalItem {
  id: number;
  nota_id: number;
  produto_id: number;
  descricao_produto: string;
  quantidade: number;
}

export interface NotaFiscal {
  id: number;
  numero_sequencial: number;
  status: 'Aberta' | 'Fechada';
  itens: NotaFiscalItem[];
}
