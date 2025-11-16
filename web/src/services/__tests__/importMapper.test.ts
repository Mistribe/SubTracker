import {describe, expect, it} from 'vitest';
import {LabelFieldMapper, ProviderFieldMapper, SubscriptionFieldMapper} from '../importMapper';
import type {DtoCreateLabelRequest, DtoCreateProviderRequest, DtoCreateSubscriptionRequest,} from '../../api';

describe('LabelFieldMapper', () => {
    const mapper = new LabelFieldMapper();

    describe('mapFields', () => {
        it('should map basic label fields correctly', () => {
            const rawRecord = {
                name: 'Entertainment',
                color: '#FF5733',
            };

            const result = mapper.mapFields(rawRecord);

            expect(result.name).toBe('Entertainment');
            expect(result.color).toBe('#FF5733');
            expect(result.owner).toEqual({type: 'personal'});
        });

        it('should map valid UUID in id field', () => {
            const rawRecord = {
                id: '550e8400-e29b-41d4-a716-446655440001',
                name: 'Entertainment',
                color: '#FF5733',
            };

            const result = mapper.mapFields(rawRecord) as any;

            expect(result.id).toBe('550e8400-e29b-41d4-a716-446655440001');
            expect(result.name).toBe('Entertainment');
        });

        it('should trim whitespace from UUID', () => {
            const rawRecord = {
                id: '  550e8400-e29b-41d4-a716-446655440001  ',
                name: 'Entertainment',
                color: '#FF5733',
            };

            const result = mapper.mapFields(rawRecord) as any;

            expect(result.id).toBe('550e8400-e29b-41d4-a716-446655440001');
        });

        it('should not map invalid UUID in id field', () => {
            const rawRecord = {
                id: 'invalid-uuid',
                name: 'Entertainment',
                color: '#FF5733',
            };

            const result = mapper.mapFields(rawRecord) as any;

            expect(result.id).toBeUndefined();
            expect(result.name).toBe('Entertainment');
        });

        it('should not map id field when empty', () => {
            const rawRecord = {
                id: '',
                name: 'Entertainment',
                color: '#FF5733',
            };

            const result = mapper.mapFields(rawRecord) as any;

            expect(result.id).toBeUndefined();
        });

        it('should not map id field when null', () => {
            const rawRecord = {
                id: null,
                name: 'Entertainment',
                color: '#FF5733',
            };

            const result = mapper.mapFields(rawRecord) as any;

            expect(result.id).toBeUndefined();
        });

        it('should add # prefix to color if missing', () => {
            const rawRecord = {
                name: 'Test',
                color: 'FF5733',
            };

            const result = mapper.mapFields(rawRecord);

            expect(result.color).toBe('#FF5733');
        });

        it('should trim whitespace from fields', () => {
            const rawRecord = {
                name: '  Entertainment  ',
                color: '  #FF5733  ',
            };

            const result = mapper.mapFields(rawRecord);

            expect(result.name).toBe('Entertainment');
            expect(result.color).toBe('#FF5733');
        });

        it('should map owner object correctly', () => {
            const rawRecord = {
                name: 'Test',
                color: '#FF5733',
                owner: {
                    type: 'family',
                    familyId: 'family-123',
                },
            };

            const result = mapper.mapFields(rawRecord);

            expect(result.owner).toEqual({
                type: 'family',
                familyId: 'family-123',
            });
        });

        it('should map owner from separate fields', () => {
            const rawRecord = {
                name: 'Test',
                color: '#FF5733',
                ownerType: 'family',
                ownerFamilyId: 'family-456',
            };

            const result = mapper.mapFields(rawRecord);

            expect(result.owner).toEqual({
                type: 'family',
                familyId: 'family-456',
            });
        });

        it('should default to personal owner if not provided', () => {
            const rawRecord = {
                name: 'Test',
                color: '#FF5733',
            };

            const result = mapper.mapFields(rawRecord);

            expect(result.owner).toEqual({type: 'personal'});
        });
    });

    describe('validate', () => {
        it('should validate a valid label', () => {
            const record: Partial<DtoCreateLabelRequest> = {
                name: 'Entertainment',
                color: '#FF5733',
                owner: 'personal',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(true);
            expect(result.errors).toHaveLength(0);
        });

        it('should validate a label with valid UUID', () => {
            const record: any = {
                id: '550e8400-e29b-41d4-a716-446655440001',
                name: 'Entertainment',
                color: '#FF5733',
                owner: 'personal',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(true);
            expect(result.errors).toHaveLength(0);
        });

        it('should reject invalid UUID format', () => {
            const record: any = {
                id: 'invalid-uuid-format',
                name: 'Entertainment',
                color: '#FF5733',
                owner: 'personal',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(false);
            expect(result.errors.some(e => e.field === 'id')).toBe(true);
            expect(result.errors.find(e => e.field === 'id')?.message).toContain('Invalid UUID format');
        });

        it('should accept empty UUID (backend will generate)', () => {
            const record: any = {
                id: '',
                name: 'Entertainment',
                color: '#FF5733',
                owner: 'personal',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(true);
            expect(result.errors).toHaveLength(0);
        });

        it('should accept null UUID (backend will generate)', () => {
            const record: any = {
                id: null,
                name: 'Entertainment',
                color: '#FF5733',
                owner: 'personal',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(true);
            expect(result.errors).toHaveLength(0);
        });

        it('should accept undefined UUID (backend will generate)', () => {
            const record: Partial<DtoCreateLabelRequest> = {
                name: 'Entertainment',
                color: '#FF5733',
                owner: 'personal',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(true);
            expect(result.errors).toHaveLength(0);
        });

        it('should validate 8-character hex colors (with alpha)', () => {
            const record: Partial<DtoCreateLabelRequest> = {
                name: 'Test',
                color: '#AAFF5733',
                owner: 'personal',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(true);
            expect(result.errors).toHaveLength(0);
        });

        it('should reject missing name', () => {
            const record: Partial<DtoCreateLabelRequest> = {
                color: '#FF5733',
                owner: 'personal',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(false);
            expect(result.errors).toContainEqual({
                field: 'name',
                message: 'Name is required',
                severity: 'error',
            });
        });

        it('should reject missing color', () => {
            const record: Partial<DtoCreateLabelRequest> = {
                name: 'Test',
                owner: 'personal',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(false);
            expect(result.errors).toContainEqual({
                field: 'color',
                message: 'Color is required',
                severity: 'error',
            });
        });

        it('should reject invalid hex color', () => {
            const record: Partial<DtoCreateLabelRequest> = {
                name: 'Test',
                color: 'not-a-color',
                owner: 'personal',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(false);
            expect(result.errors).toContainEqual({
                field: 'color',
                message: 'Color must be a valid hex color (e.g., #FF5733 or #AAFF5733)',
                severity: 'error',
            });
        });

        it('should reject missing owner', () => {
            const record: Partial<DtoCreateLabelRequest> = {
                name: 'Test',
                color: '#FF5733',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(false);
            expect(result.errors).toContainEqual({
                field: 'owner',
                message: 'Owner is required',
                severity: 'error',
            });
        });

        it('should reject invalid owner type', () => {
            const record: Partial<DtoCreateLabelRequest> = {
                name: 'Test',
                color: '#FF5733',
                owner: 'invalid' as any,
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(false);
            expect(result.errors).toContainEqual({
                field: 'owner.type',
                message: 'Owner type must be one of: personal, family, system',
                severity: 'error',
            });
        });

        it('should reject family owner without familyId', () => {
            const record: Partial<DtoCreateLabelRequest> = {
                name: 'Test',
                color: '#FF5733',
                owner: 'family',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(false);
            expect(result.errors).toContainEqual({
                field: 'owner.familyId',
                message: 'Family ID is required when owner type is family',
                severity: 'error',
            });
        });
    });
});

describe('ProviderFieldMapper', () => {
    const mapper = new ProviderFieldMapper();

    describe('mapFields', () => {
        it('should map basic provider fields correctly', () => {
            const rawRecord = {
                name: 'Netflix',
                description: 'Streaming service',
                url: 'https://netflix.com',
            };

            const result = mapper.mapFields(rawRecord);

            expect(result.name).toBe('Netflix');
            expect(result.description).toBe('Streaming service');
            expect(result.url).toBe('https://netflix.com');
        });

        it('should map valid UUID in id field', () => {
            const rawRecord = {
                id: '550e8400-e29b-41d4-a716-446655440010',
                name: 'Netflix',
            };

            const result = mapper.mapFields(rawRecord) as any;

            expect(result.id).toBe('550e8400-e29b-41d4-a716-446655440010');
            expect(result.name).toBe('Netflix');
        });

        it('should trim whitespace from UUID', () => {
            const rawRecord = {
                id: '  550e8400-e29b-41d4-a716-446655440010  ',
                name: 'Netflix',
            };

            const result = mapper.mapFields(rawRecord) as any;

            expect(result.id).toBe('550e8400-e29b-41d4-a716-446655440010');
        });

        it('should not map invalid UUID in id field', () => {
            const rawRecord = {
                id: 'not-a-uuid',
                name: 'Netflix',
            };

            const result = mapper.mapFields(rawRecord) as any;

            expect(result.id).toBeUndefined();
        });

        it('should not map id field when empty', () => {
            const rawRecord = {
                id: '',
                name: 'Netflix',
            };

            const result = mapper.mapFields(rawRecord) as any;

            expect(result.id).toBeUndefined();
        });

        it('should trim whitespace from fields', () => {
            const rawRecord = {
                name: '  Netflix  ',
                description: '  Streaming  ',
            };

            const result = mapper.mapFields(rawRecord);

            expect(result.name).toBe('Netflix');
            expect(result.description).toBe('Streaming');
        });

        it('should map all optional URL fields', () => {
            const rawRecord = {
                name: 'AWS',
                url: 'https://aws.amazon.com',
                iconUrl: 'https://aws.amazon.com/icon.png',
                pricingPageUrl: 'https://aws.amazon.com/pricing',
            };

            const result = mapper.mapFields(rawRecord);

            expect(result.url).toBe('https://aws.amazon.com');
            expect(result.iconUrl).toBe('https://aws.amazon.com/icon.png');
            expect(result.pricingPageUrl).toBe('https://aws.amazon.com/pricing');
        });

        it('should map labels from array', () => {
            const rawRecord = {
                name: 'Netflix',
                labels: ['entertainment', 'streaming'],
            };

            const result = mapper.mapFields(rawRecord);

            expect(result.labels).toEqual(['entertainment', 'streaming']);
        });

        it('should map labels from comma-separated string', () => {
            const rawRecord = {
                name: 'Netflix',
                labels: 'entertainment, streaming, video',
            };

            const result = mapper.mapFields(rawRecord);

            expect(result.labels).toEqual(['entertainment', 'streaming', 'video']);
        });

        it('should filter empty labels', () => {
            const rawRecord = {
                name: 'Netflix',
                labels: ['entertainment', '', '  ', 'streaming'],
            };

            const result = mapper.mapFields(rawRecord);

            expect(result.labels).toEqual(['entertainment', 'streaming']);
        });

        it('should map owner if provided', () => {
            const rawRecord = {
                name: 'Netflix',
                owner: {
                    type: 'family',
                    familyId: 'family-123',
                },
            };

            const result = mapper.mapFields(rawRecord);

            expect(result.owner).toEqual({
                type: 'family',
                familyId: 'family-123',
            });
        });

        it('should not set owner if not provided', () => {
            const rawRecord = {
                name: 'Netflix',
            };

            const result = mapper.mapFields(rawRecord);

            expect(result.owner).toBeUndefined();
        });
    });

    describe('validate', () => {
        it('should validate a valid provider with only required fields', () => {
            const record: Partial<DtoCreateProviderRequest> = {
                name: 'Netflix',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(true);
            expect(result.errors).toHaveLength(0);
        });

        it('should validate a provider with valid UUID', () => {
            const record: any = {
                id: '550e8400-e29b-41d4-a716-446655440010',
                name: 'Netflix',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(true);
            expect(result.errors).toHaveLength(0);
        });

        it('should reject invalid UUID format', () => {
            const record: any = {
                id: 'bad-uuid',
                name: 'Netflix',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(false);
            expect(result.errors.some(e => e.field === 'id')).toBe(true);
            expect(result.errors.find(e => e.field === 'id')?.message).toContain('Invalid UUID format');
        });

        it('should accept empty UUID', () => {
            const record: any = {
                id: '',
                name: 'Netflix',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(true);
            expect(result.errors).toHaveLength(0);
        });

        it('should validate a provider with all fields', () => {
            const record: Partial<DtoCreateProviderRequest> = {
                name: 'Netflix',
                description: 'Streaming',
                url: 'https://netflix.com',
                iconUrl: 'https://netflix.com/icon.png',
                pricingPageUrl: 'https://netflix.com/pricing',
                labels: ['entertainment'],
                owner: 'personal',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(true);
            expect(result.errors).toHaveLength(0);
        });

        it('should reject missing name', () => {
            const record: Partial<DtoCreateProviderRequest> = {
                description: 'Test',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(false);
            expect(result.errors).toContainEqual({
                field: 'name',
                message: 'Name is required',
                severity: 'error',
            });
        });

        it('should reject invalid URL', () => {
            const record: Partial<DtoCreateProviderRequest> = {
                name: 'Netflix',
                url: 'not-a-url',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(false);
            expect(result.errors).toContainEqual({
                field: 'url',
                message: 'URL must be a valid URL',
                severity: 'error',
            });
        });

        it('should reject invalid iconUrl', () => {
            const record: Partial<DtoCreateProviderRequest> = {
                name: 'Netflix',
                iconUrl: 'invalid-url',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(false);
            expect(result.errors).toContainEqual({
                field: 'iconUrl',
                message: 'Icon URL must be a valid URL',
                severity: 'error',
            });
        });

        it('should reject invalid pricingPageUrl', () => {
            const record: Partial<DtoCreateProviderRequest> = {
                name: 'Netflix',
                pricingPageUrl: 'bad-url',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(false);
            expect(result.errors).toContainEqual({
                field: 'pricingPageUrl',
                message: 'Pricing page URL must be a valid URL',
                severity: 'error',
            });
        });

        it('should validate owner if provided', () => {
            const record: Partial<DtoCreateProviderRequest> = {
                name: 'Netflix',
                owner: 'family',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(false);
            expect(result.errors).toContainEqual({
                field: 'owner.familyId',
                message: 'Family ID is required when owner type is family',
                severity: 'error',
            });
        });
    });
});

describe('SubscriptionFieldMapper', () => {
    const mapper = new SubscriptionFieldMapper();

    describe('mapFields', () => {
        it('should map basic subscription fields correctly', () => {
            const rawRecord = {
                providerId: 'provider-123',
                startDate: '2024-01-01',
                recurrency: 'monthly',
            };

            const result = mapper.mapFields(rawRecord);

            expect(result.providerId).toBe('provider-123');
            expect(result.startDate).toEqual(new Date('2024-01-01'));
            expect(result.recurrency).toBe('monthly');
            expect(result.owner).toEqual({type: 'personal'});
        });

        it('should map valid UUID in id field', () => {
            const rawRecord = {
                id: '550e8400-e29b-41d4-a716-446655440020',
                providerId: 'provider-123',
                startDate: '2024-01-01',
                recurrency: 'monthly',
            };

            const result = mapper.mapFields(rawRecord) as any;

            expect(result.id).toBe('550e8400-e29b-41d4-a716-446655440020');
            expect(result.providerId).toBe('provider-123');
        });

        it('should trim whitespace from UUID', () => {
            const rawRecord = {
                id: '  550e8400-e29b-41d4-a716-446655440020  ',
                providerId: 'provider-123',
                startDate: '2024-01-01',
                recurrency: 'monthly',
            };

            const result = mapper.mapFields(rawRecord) as any;

            expect(result.id).toBe('550e8400-e29b-41d4-a716-446655440020');
        });

        it('should not map invalid UUID in id field', () => {
            const rawRecord = {
                id: 'invalid',
                providerId: 'provider-123',
                startDate: '2024-01-01',
                recurrency: 'monthly',
            };

            const result = mapper.mapFields(rawRecord) as any;

            expect(result.id).toBeUndefined();
        });

        it('should not map id field when empty', () => {
            const rawRecord = {
                id: '',
                providerId: 'provider-123',
                startDate: '2024-01-01',
                recurrency: 'monthly',
            };

            const result = mapper.mapFields(rawRecord) as any;

            expect(result.id).toBeUndefined();
        });

        it('should normalize recurrency to lowercase', () => {
            const rawRecord = {
                providerId: 'provider-123',
                startDate: '2024-01-01',
                recurrency: 'MONTHLY',
            };

            const result = mapper.mapFields(rawRecord);

            expect(result.recurrency).toBe('monthly');
        });

        it('should map optional fields', () => {
            const rawRecord = {
                providerId: 'provider-123',
                startDate: '2024-01-01',
                recurrency: 'monthly',
                friendlyName: 'Netflix Premium',
                endDate: '2024-12-31',
                customRecurrency: 3,
            };

            const result = mapper.mapFields(rawRecord);

            expect(result.friendlyName).toBe('Netflix Premium');
            expect(result.endDate).toEqual(new Date('2024-12-31'));
            expect(result.customRecurrency).toBe(3);
        });

        it('should map custom price from object', () => {
            const rawRecord = {
                providerId: 'provider-123',
                startDate: '2024-01-01',
                recurrency: 'monthly',
                customPrice: {
                    value: 15.99,
                    currency: 'usd',
                },
            };

            const result = mapper.mapFields(rawRecord);

            expect(result.customPrice).toEqual({
                value: 15.99,
                currency: 'USD',
            });
        });

        it('should map custom price from separate fields', () => {
            const rawRecord = {
                providerId: 'provider-123',
                startDate: '2024-01-01',
                recurrency: 'monthly',
                customPriceAmount: 15.99,
                customPriceCurrency: 'eur',
            };

            const result = mapper.mapFields(rawRecord);

            expect(result.customPrice).toEqual({
                value: 15.99,
                currency: 'EUR',
            });
        });

        it('should map labels from array', () => {
            const rawRecord = {
                providerId: 'provider-123',
                startDate: '2024-01-01',
                recurrency: 'monthly',
                labels: ['entertainment', 'streaming'],
            };

            const result = mapper.mapFields(rawRecord);

            expect(result.labels).toEqual(['entertainment', 'streaming']);
        });

        it('should map labels from comma-separated string', () => {
            const rawRecord = {
                providerId: 'provider-123',
                startDate: '2024-01-01',
                recurrency: 'monthly',
                labels: 'entertainment, streaming',
            };

            const result = mapper.mapFields(rawRecord);

            expect(result.labels).toEqual(['entertainment', 'streaming']);
        });

        it('should map payer from object', () => {
            const rawRecord = {
                providerId: 'provider-123',
                startDate: '2024-01-01',
                recurrency: 'monthly',
                payer: {
                    type: 'family',
                    familyId: 'family-123',
                },
            };

            const result = mapper.mapFields(rawRecord);

            expect(result.payer).toEqual({
                type: 'family',
                familyId: 'family-123',
            });
        });

        it('should map payer from separate fields', () => {
            const rawRecord = {
                providerId: 'provider-123',
                startDate: '2024-01-01',
                recurrency: 'monthly',
                payerType: 'family_member',
                payerFamilyId: 'family-123',
                payerMemberId: 'member-456',
            };

            const result = mapper.mapFields(rawRecord);

            expect(result.payer).toEqual({
                type: 'family_member',
                familyId: 'family-123',
                memberId: 'member-456',
            });
        });

        it('should map free trial from object', () => {
            const rawRecord = {
                providerId: 'provider-123',
                startDate: '2024-01-01',
                recurrency: 'monthly',
                freeTrial: {
                    startDate: '2024-01-01',
                    endDate: '2024-01-31',
                },
            };

            const result = mapper.mapFields(rawRecord);

            expect(result.freeTrial).toEqual({
                startDate: new Date('2024-01-01'),
                endDate: new Date('2024-01-31'),
            });
        });

        it('should map free trial from separate fields', () => {
            const rawRecord = {
                providerId: 'provider-123',
                startDate: '2024-01-01',
                recurrency: 'monthly',
                freeTrialStartDate: '2024-01-01',
                freeTrialEndDate: '2024-01-31',
            };

            const result = mapper.mapFields(rawRecord);

            expect(result.freeTrial).toEqual({
                startDate: new Date('2024-01-01'),
                endDate: new Date('2024-01-31'),
            });
        });

        it('should map family users from array', () => {
            const rawRecord = {
                providerId: 'provider-123',
                startDate: '2024-01-01',
                recurrency: 'monthly',
                familyUsers: ['user-1', 'user-2'],
            };

            const result = mapper.mapFields(rawRecord);

            expect(result.familyUsers).toEqual(['user-1', 'user-2']);
        });

        it('should map family users from comma-separated string', () => {
            const rawRecord = {
                providerId: 'provider-123',
                startDate: '2024-01-01',
                recurrency: 'monthly',
                familyUsers: 'user-1, user-2, user-3',
            };

            const result = mapper.mapFields(rawRecord);

            expect(result.familyUsers).toEqual(['user-1', 'user-2', 'user-3']);
        });
    });

    describe('validate', () => {
        it('should validate a valid subscription with required fields', () => {
            const record: Partial<DtoCreateSubscriptionRequest> = {
                providerId: 'provider-123',
                startDate: new Date('2024-01-01'),
                recurrency: 'monthly',
                owner: 'personal',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(true);
            expect(result.errors).toHaveLength(0);
        });

        it('should validate a subscription with valid UUID', () => {
            const record: any = {
                id: '550e8400-e29b-41d4-a716-446655440020',
                providerId: 'provider-123',
                startDate: new Date('2024-01-01'),
                recurrency: 'monthly',
                owner: 'personal',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(true);
            expect(result.errors).toHaveLength(0);
        });

        it('should reject invalid UUID format', () => {
            const record: any = {
                id: 'not-valid-uuid',
                providerId: 'provider-123',
                startDate: new Date('2024-01-01'),
                recurrency: 'monthly',
                owner: 'personal',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(false);
            expect(result.errors.some(e => e.field === 'id')).toBe(true);
            expect(result.errors.find(e => e.field === 'id')?.message).toContain('Invalid UUID format');
        });

        it('should accept empty UUID', () => {
            const record: any = {
                id: '',
                providerId: 'provider-123',
                startDate: new Date('2024-01-01'),
                recurrency: 'monthly',
                owner: 'personal',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(true);
            expect(result.errors).toHaveLength(0);
        });

        it('should accept null UUID', () => {
            const record: any = {
                id: null,
                providerId: 'provider-123',
                startDate: new Date('2024-01-01'),
                recurrency: 'monthly',
                owner: 'personal',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(true);
            expect(result.errors).toHaveLength(0);
        });

        it('should validate all recurrency types', () => {
            const recurrencies = ['daily', 'weekly', 'monthly', 'quarterly', 'yearly', 'custom'];

            recurrencies.forEach((recurrency) => {
                const record: Partial<DtoCreateSubscriptionRequest> = {
                    providerId: 'provider-123',
                    startDate: new Date('2024-01-01'),
                    recurrency,
                    owner: 'personal',
                };

                const result = mapper.validate(record);

                expect(result.isValid).toBe(true);
            });
        });

        it('should reject missing providerId', () => {
            const record: Partial<DtoCreateSubscriptionRequest> = {
                startDate: new Date('2024-01-01'),
                recurrency: 'monthly',
                owner: 'personal',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(false);
            expect(result.errors).toContainEqual({
                field: 'providerId',
                message: 'Provider ID is required',
                severity: 'error',
            });
        });

        it('should reject missing startDate', () => {
            const record: Partial<DtoCreateSubscriptionRequest> = {
                providerId: 'provider-123',
                recurrency: 'monthly',
                owner: 'personal',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(false);
            expect(result.errors).toContainEqual({
                field: 'startDate',
                message: 'Start date is required',
                severity: 'error',
            });
        });

        it('should reject invalid startDate', () => {
            const record: Partial<DtoCreateSubscriptionRequest> = {
                providerId: 'provider-123',
                startDate: new Date('invalid'),
                recurrency: 'monthly',
                owner: 'personal',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(false);
            expect(result.errors).toContainEqual({
                field: 'startDate',
                message: 'Start date must be a valid date',
                severity: 'error',
            });
        });

        it('should reject missing recurrency', () => {
            const record: Partial<DtoCreateSubscriptionRequest> = {
                providerId: 'provider-123',
                startDate: new Date('2024-01-01'),
                owner: 'personal',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(false);
            expect(result.errors).toContainEqual({
                field: 'recurrency',
                message: 'Recurrency is required',
                severity: 'error',
            });
        });

        it('should reject invalid recurrency', () => {
            const record: Partial<DtoCreateSubscriptionRequest> = {
                providerId: 'provider-123',
                startDate: new Date('2024-01-01'),
                recurrency: 'invalid',
                owner: 'personal',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(false);
            expect(result.errors.some((e) => e.field === 'recurrency')).toBe(true);
        });

        it('should reject missing owner', () => {
            const record: Partial<DtoCreateSubscriptionRequest> = {
                providerId: 'provider-123',
                startDate: new Date('2024-01-01'),
                recurrency: 'monthly',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(false);
            expect(result.errors).toContainEqual({
                field: 'owner',
                message: 'Owner is required',
                severity: 'error',
            });
        });

        it('should reject endDate before startDate', () => {
            const record: Partial<DtoCreateSubscriptionRequest> = {
                providerId: 'provider-123',
                startDate: new Date('2024-12-31'),
                endDate: new Date('2024-01-01'),
                recurrency: 'monthly',
                owner: 'personal',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(false);
            expect(result.errors).toContainEqual({
                field: 'endDate',
                message: 'End date must be after start date',
                severity: 'error',
            });
        });

        it('should reject invalid customRecurrency', () => {
            const record: Partial<DtoCreateSubscriptionRequest> = {
                providerId: 'provider-123',
                startDate: new Date('2024-01-01'),
                recurrency: 'custom',
                customRecurrency: -1,
                owner: 'personal',
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(false);
            expect(result.errors).toContainEqual({
                field: 'customRecurrency',
                message: 'Custom recurrency must be a positive number',
                severity: 'error',
            });
        });

        it('should validate custom price', () => {
            const record: Partial<DtoCreateSubscriptionRequest> = {
                providerId: 'provider-123',
                startDate: new Date('2024-01-01'),
                recurrency: 'monthly',
                owner: 'personal',
                customPrice: {
                    value: 15.99,
                    currency: 'USD',
                },
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(true);
            expect(result.errors).toHaveLength(0);
        });

        it('should reject invalid custom price value', () => {
            const record: Partial<DtoCreateSubscriptionRequest> = {
                providerId: 'provider-123',
                startDate: new Date('2024-01-01'),
                recurrency: 'monthly',
                owner: 'personal',
                customPrice: {
                    value: -10,
                    currency: 'USD',
                },
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(false);
            expect(result.errors).toContainEqual({
                field: 'customPrice.value',
                message: 'Price value must be a non-negative number',
                severity: 'error',
            });
        });

        it('should reject invalid currency code', () => {
            const record: Partial<DtoCreateSubscriptionRequest> = {
                providerId: 'provider-123',
                startDate: new Date('2024-01-01'),
                recurrency: 'monthly',
                owner: 'personal',
                customPrice: {
                    value: 15.99,
                    currency: 'US',
                },
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(false);
            expect(result.errors).toContainEqual({
                field: 'customPrice.currency',
                message: 'Currency must be a 3-letter ISO code (e.g., USD, EUR)',
                severity: 'error',
            });
        });

        it('should validate payer', () => {
            const record: Partial<DtoCreateSubscriptionRequest> = {
                providerId: 'provider-123',
                startDate: new Date('2024-01-01'),
                recurrency: 'monthly',
                owner: 'personal',
                payer: {
                    type: 'family',
                },
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(true);
            expect(result.errors).toHaveLength(0);
        });

        it('should reject payer without familyId', () => {
            const record: Partial<DtoCreateSubscriptionRequest> = {
                providerId: 'provider-123',
                startDate: new Date('2024-01-01'),
                recurrency: 'monthly',
                owner: 'personal',
                payer: {
                    type: 'family',
                },
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(false);
            expect(result.errors).toContainEqual({
                field: 'payer.familyId',
                message: 'Payer family ID is required',
                severity: 'error',
            });
        });

        it('should reject family_member payer without memberId', () => {
            const record: Partial<DtoCreateSubscriptionRequest> = {
                providerId: 'provider-123',
                startDate: new Date('2024-01-01'),
                recurrency: 'monthly',
                owner: 'personal',
                payer: {
                    type: 'family_member',
                },
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(false);
            expect(result.errors).toContainEqual({
                field: 'payer.memberId',
                message: 'Member ID is required when payer type is family_member',
                severity: 'error',
            });
        });

        it('should validate free trial', () => {
            const record: Partial<DtoCreateSubscriptionRequest> = {
                providerId: 'provider-123',
                startDate: new Date('2024-01-01'),
                recurrency: 'monthly',
                owner: 'personal',
                freeTrial: {
                    startDate: new Date('2024-01-01'),
                    endDate: new Date('2024-01-31'),
                },
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(true);
            expect(result.errors).toHaveLength(0);
        });

        it('should reject free trial with endDate before startDate', () => {
            const record: Partial<DtoCreateSubscriptionRequest> = {
                providerId: 'provider-123',
                startDate: new Date('2024-01-01'),
                recurrency: 'monthly',
                owner: 'personal',
                freeTrial: {
                    startDate: new Date('2024-01-31'),
                    endDate: new Date('2024-01-01'),
                },
            };

            const result = mapper.validate(record);

            expect(result.isValid).toBe(false);
            expect(result.errors).toContainEqual({
                field: 'freeTrial.endDate',
                message: 'Free trial end date must be after start date',
                severity: 'error',
            });
        });
    });
});

describe('Mixed ID scenarios', () => {
    describe('LabelFieldMapper with mixed IDs', () => {
        const mapper = new LabelFieldMapper();

        it('should handle batch with some records having IDs and some without', () => {
            const records = [
                {
                    id: '550e8400-e29b-41d4-a716-446655440001',
                    name: 'Label 1',
                    color: '#FF5733',
                },
                {
                    name: 'Label 2',
                    color: '#33FF57',
                },
                {
                    id: '550e8400-e29b-41d4-a716-446655440003',
                    name: 'Label 3',
                    color: '#3357FF',
                },
            ];

            const results = records.map(r => mapper.mapFields(r));

            expect((results[0] as any).id).toBe('550e8400-e29b-41d4-a716-446655440001');
            expect((results[1] as any).id).toBeUndefined();
            expect((results[2] as any).id).toBe('550e8400-e29b-41d4-a716-446655440003');
        });

        it('should validate mixed batch correctly', () => {
            const records: any[] = [
                {
                    id: '550e8400-e29b-41d4-a716-446655440001',
                    name: 'Label 1',
                    color: '#FF5733',
                    owner: 'personal',
                },
                {
                    name: 'Label 2',
                    color: '#33FF57',
                    owner: 'personal',
                },
                {
                    id: 'invalid-uuid',
                    name: 'Label 3',
                    color: '#3357FF',
                    owner: 'personal',
                },
            ];

            const results = records.map(r => mapper.validate(r));

            expect(results[0].isValid).toBe(true);
            expect(results[1].isValid).toBe(true);
            expect(results[2].isValid).toBe(false);
            expect(results[2].errors.some(e => e.field === 'id')).toBe(true);
        });
    });

    describe('ProviderFieldMapper with mixed IDs', () => {
        const mapper = new ProviderFieldMapper();

        it('should handle batch with some records having IDs and some without', () => {
            const records = [
                {
                    id: '550e8400-e29b-41d4-a716-446655440010',
                    name: 'Provider 1',
                },
                {
                    name: 'Provider 2',
                },
                {
                    id: '',
                    name: 'Provider 3',
                },
            ];

            const results = records.map(r => mapper.mapFields(r));

            expect((results[0] as any).id).toBe('550e8400-e29b-41d4-a716-446655440010');
            expect((results[1] as any).id).toBeUndefined();
            expect((results[2] as any).id).toBeUndefined();
        });
    });

    describe('SubscriptionFieldMapper with mixed IDs', () => {
        const mapper = new SubscriptionFieldMapper();

        it('should handle batch with some records having IDs and some without', () => {
            const records = [
                {
                    id: '550e8400-e29b-41d4-a716-446655440020',
                    providerId: 'provider-123',
                    startDate: '2024-01-01',
                    recurrency: 'monthly',
                },
                {
                    providerId: 'provider-123',
                    startDate: '2024-01-01',
                    recurrency: 'monthly',
                },
                {
                    id: null,
                    providerId: 'provider-123',
                    startDate: '2024-01-01',
                    recurrency: 'monthly',
                },
            ];

            const results = records.map(r => mapper.mapFields(r));

            expect((results[0] as any).id).toBe('550e8400-e29b-41d4-a716-446655440020');
            expect((results[1] as any).id).toBeUndefined();
            expect((results[2] as any).id).toBeUndefined();
        });
    });
});
