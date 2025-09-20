import type { CurrenciesApi } from '@/api/apis/CurrenciesApi';
import type { FamilyApi } from '@/api/apis/FamilyApi';
import type { LabelsApi } from '@/api/apis/LabelsApi';
import type { ProvidersApi } from '@/api/apis/ProvidersApi';
import type { SubscriptionsApi } from '@/api/apis/SubscriptionsApi';
import type { UsersApi } from '@/api/apis/UsersApi';

export interface ApiClient {
  currencies: CurrenciesApi;
  families: FamilyApi;
  labels: LabelsApi;
  providers: ProvidersApi;
  subscriptions: SubscriptionsApi;
  users: UsersApi;
}
