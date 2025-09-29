import type { CurrenciesApi } from '@/api/apis/CurrenciesApi';
import type { FamilyApi } from '@/api/apis/FamilyApi';
import type { LabelsApi } from '@/api/apis/LabelsApi';
import type { ProvidersApi } from '@/api/apis/ProvidersApi';
import type { SubscriptionsApi } from '@/api/apis/SubscriptionsApi';
import type { AccountsApi } from '@/api/apis/AccountsApi';

export interface ApiClient {
  currencies: CurrenciesApi;
  families: FamilyApi;
  labels: LabelsApi;
  providers: ProvidersApi;
  subscriptions: SubscriptionsApi;
  accounts: AccountsApi;
}
