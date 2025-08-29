import type { AuthenticationProvider } from '@microsoft/kiota-abstractions';
import type { RequestInformation } from '@microsoft/kiota-abstractions/dist/es/src/requestInformation';

// Generic token auth provider compatible with Kiota's FetchRequestAdapter
export class AuthTokenProvider implements AuthenticationProvider {
  private readonly getToken: () => Promise<string | null | undefined>;

  constructor(getToken: () => Promise<string | null | undefined>) {
    this.getToken = getToken;
  }

  async authenticateRequest(request: RequestInformation): Promise<void> {
    const token = await this.getToken();
    if (token) {
      request.headers.add('Authorization', `Bearer ${token}`);
    }
  }
}
