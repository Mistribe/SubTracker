import type {AuthenticationProvider} from '@microsoft/kiota-abstractions';
import type {RequestInformation} from "@microsoft/kiota-abstractions/dist/es/src/requestInformation";

export class KindeAuthProvider implements AuthenticationProvider {
    private readonly getToken: () => Promise<string | undefined>;

    constructor(getToken: () => Promise<string | undefined>) {
        this.getToken = getToken;
    }

    async authenticateRequest(request: RequestInformation): Promise<void> {
        const token = await this.getToken();
        if (token) {
            request.headers.add('Authorization', `Bearer ${token}`);
        }
    }
}