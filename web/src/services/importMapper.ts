import type {
  DtoCreateLabelRequest,
  DtoCreateProviderRequest,
  DtoCreateSubscriptionRequest,
  DtoEditableOwnerModel,
  DtoEditableOwnerModelTypeEnum,
  DtoAmountModel,
  DtoEditableSubscriptionPayerModel,
  DtoEditableSubscriptionPayerModelTypeEnum,
  DtoSubscriptionFreeTrialModel,
} from '../api';
import type { FieldMapper, ValidationResult, ValidationError } from '../types/import';

/**
 * Base field mapper implementation with common validation utilities
 */
abstract class BaseFieldMapper<T> implements FieldMapper<T> {
  abstract mapFields(rawRecord: Record<string, any>): Partial<T>;
  abstract validate(record: Partial<T>): ValidationResult;

  /**
   * Validates if a string is a valid hex color
   */
  protected isValidHexColor(color: string): boolean {
    return /^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{8})$/.test(color);
  }

  /**
   * Validates if a string is a valid URL
   */
  protected isValidUrl(url: string): boolean {
    try {
      new URL(url);
      return true;
    } catch {
      return false;
    }
  }

  /**
   * Validates if a string is a valid date
   */
  protected isValidDate(dateString: string): boolean {
    const date = new Date(dateString);
    return !isNaN(date.getTime());
  }

  /**
   * Parses a date string to Date object
   */
  protected parseDate(dateString: string): Date | undefined {
    if (!dateString) return undefined;
    const date = new Date(dateString);
    return isNaN(date.getTime()) ? undefined : date;
  }

  /**
   * Creates a validation error
   */
  protected createError(field: string, message: string, severity: 'error' | 'warning' = 'error'): ValidationError {
    return { field, message, severity };
  }

  /**
   * Checks if a value is empty (null, undefined, or empty string)
   */
  protected isEmpty(value: any): boolean {
    return value === null || value === undefined || value === '';
  }
}

/**
 * Field mapper for Label entities
 */
export class LabelFieldMapper extends BaseFieldMapper<DtoCreateLabelRequest> {
  mapFields(rawRecord: Record<string, any>): Partial<DtoCreateLabelRequest> {
    const mapped: Partial<DtoCreateLabelRequest> = {};

    // Map name
    if (!this.isEmpty(rawRecord.name)) {
      mapped.name = String(rawRecord.name).trim();
    }

    // Map color
    if (!this.isEmpty(rawRecord.color)) {
      let color = String(rawRecord.color).trim();
      // Add # prefix if missing
      if (color && !color.startsWith('#')) {
        color = '#' + color;
      }
      mapped.color = color;
    }

    // Map owner (default to personal if not provided)
    if (rawRecord.owner && typeof rawRecord.owner === 'object') {
      mapped.owner = this.mapOwner(rawRecord.owner);
    } else if (!this.isEmpty(rawRecord.ownerType)) {
      mapped.owner = this.mapOwner({
        type: rawRecord.ownerType,
        familyId: rawRecord.ownerFamilyId,
      });
    } else {
      // Default to personal ownership
      mapped.owner = { type: 'personal' as DtoEditableOwnerModelTypeEnum };
    }

    return mapped;
  }

  validate(record: Partial<DtoCreateLabelRequest>): ValidationResult {
    const errors: ValidationError[] = [];

    // Validate required fields
    if (this.isEmpty(record.name)) {
      errors.push(this.createError('name', 'Name is required'));
    }

    if (this.isEmpty(record.color)) {
      errors.push(this.createError('color', 'Color is required'));
    } else if (!this.isValidHexColor(record.color!)) {
      errors.push(this.createError('color', 'Color must be a valid hex color (e.g., #FF5733 or #AAFF5733)'));
    }

    if (!record.owner) {
      errors.push(this.createError('owner', 'Owner is required'));
    } else {
      const ownerErrors = this.validateOwner(record.owner);
      errors.push(...ownerErrors);
    }

    return {
      isValid: errors.length === 0,
      errors,
    };
  }

  private mapOwner(rawOwner: Record<string, any>): DtoEditableOwnerModel {
    const owner: DtoEditableOwnerModel = {
      type: (rawOwner.type || 'personal') as DtoEditableOwnerModelTypeEnum,
    };

    if (rawOwner.familyId) {
      owner.familyId = String(rawOwner.familyId);
    }

    return owner;
  }

  private validateOwner(owner: DtoEditableOwnerModel): ValidationError[] {
    const errors: ValidationError[] = [];

    if (!owner.type) {
      errors.push(this.createError('owner.type', 'Owner type is required'));
    } else if (!['personal', 'family', 'system'].includes(owner.type)) {
      errors.push(this.createError('owner.type', 'Owner type must be one of: personal, family, system'));
    }

    if (owner.type === 'family' && this.isEmpty(owner.familyId)) {
      errors.push(this.createError('owner.familyId', 'Family ID is required when owner type is family'));
    }

    return errors;
  }
}

/**
 * Field mapper for Provider entities
 */
export class ProviderFieldMapper extends BaseFieldMapper<DtoCreateProviderRequest> {
  mapFields(rawRecord: Record<string, any>): Partial<DtoCreateProviderRequest> {
    const mapped: Partial<DtoCreateProviderRequest> = {};

    // Map name (required)
    if (!this.isEmpty(rawRecord.name)) {
      mapped.name = String(rawRecord.name).trim();
    }

    // Map optional fields
    if (!this.isEmpty(rawRecord.description)) {
      mapped.description = String(rawRecord.description).trim();
    }

    if (!this.isEmpty(rawRecord.url)) {
      mapped.url = String(rawRecord.url).trim();
    }

    if (!this.isEmpty(rawRecord.iconUrl)) {
      mapped.iconUrl = String(rawRecord.iconUrl).trim();
    }

    if (!this.isEmpty(rawRecord.pricingPageUrl)) {
      mapped.pricingPageUrl = String(rawRecord.pricingPageUrl).trim();
    }

    // Map labels (array of strings)
    if (rawRecord.labels) {
      if (Array.isArray(rawRecord.labels)) {
        mapped.labels = rawRecord.labels.map((label: any) => String(label).trim()).filter(Boolean);
      } else if (typeof rawRecord.labels === 'string') {
        // Support comma-separated string
        mapped.labels = rawRecord.labels.split(',').map((label: string) => label.trim()).filter(Boolean);
      }
    }

    // Map owner (optional, defaults to personal if not provided)
    if (rawRecord.owner && typeof rawRecord.owner === 'object') {
      mapped.owner = this.mapOwner(rawRecord.owner);
    } else if (!this.isEmpty(rawRecord.ownerType)) {
      mapped.owner = this.mapOwner({
        type: rawRecord.ownerType,
        familyId: rawRecord.ownerFamilyId,
      });
    }

    return mapped;
  }

  validate(record: Partial<DtoCreateProviderRequest>): ValidationResult {
    const errors: ValidationError[] = [];

    // Validate required fields
    if (this.isEmpty(record.name)) {
      errors.push(this.createError('name', 'Name is required'));
    }

    // Validate optional URL fields
    if (record.url && !this.isValidUrl(record.url)) {
      errors.push(this.createError('url', 'URL must be a valid URL'));
    }

    if (record.iconUrl && !this.isValidUrl(record.iconUrl)) {
      errors.push(this.createError('iconUrl', 'Icon URL must be a valid URL'));
    }

    if (record.pricingPageUrl && !this.isValidUrl(record.pricingPageUrl)) {
      errors.push(this.createError('pricingPageUrl', 'Pricing page URL must be a valid URL'));
    }

    // Validate owner if provided
    if (record.owner) {
      const ownerErrors = this.validateOwner(record.owner);
      errors.push(...ownerErrors);
    }

    return {
      isValid: errors.length === 0,
      errors,
    };
  }

  private mapOwner(rawOwner: Record<string, any>): DtoEditableOwnerModel {
    const owner: DtoEditableOwnerModel = {
      type: (rawOwner.type || 'personal') as DtoEditableOwnerModelTypeEnum,
    };

    if (rawOwner.familyId) {
      owner.familyId = String(rawOwner.familyId);
    }

    return owner;
  }

  private validateOwner(owner: DtoEditableOwnerModel): ValidationError[] {
    const errors: ValidationError[] = [];

    if (!owner.type) {
      errors.push(this.createError('owner.type', 'Owner type is required'));
    } else if (!['personal', 'family', 'system'].includes(owner.type)) {
      errors.push(this.createError('owner.type', 'Owner type must be one of: personal, family, system'));
    }

    if (owner.type === 'family' && this.isEmpty(owner.familyId)) {
      errors.push(this.createError('owner.familyId', 'Family ID is required when owner type is family'));
    }

    return errors;
  }
}

/**
 * Field mapper for Subscription entities
 */
export class SubscriptionFieldMapper extends BaseFieldMapper<DtoCreateSubscriptionRequest> {
  mapFields(rawRecord: Record<string, any>): Partial<DtoCreateSubscriptionRequest> {
    const mapped: Partial<DtoCreateSubscriptionRequest> = {};

    // Map required fields
    if (!this.isEmpty(rawRecord.providerId)) {
      mapped.providerId = String(rawRecord.providerId).trim();
    }

    if (!this.isEmpty(rawRecord.startDate)) {
      mapped.startDate = this.parseDate(String(rawRecord.startDate));
    }

    if (!this.isEmpty(rawRecord.recurrency)) {
      mapped.recurrency = String(rawRecord.recurrency).trim().toLowerCase();
    }

    // Map owner (required)
    if (rawRecord.owner && typeof rawRecord.owner === 'object') {
      mapped.owner = this.mapOwner(rawRecord.owner);
    } else if (!this.isEmpty(rawRecord.ownerType)) {
      mapped.owner = this.mapOwner({
        type: rawRecord.ownerType,
        familyId: rawRecord.ownerFamilyId,
      });
    } else {
      // Default to personal ownership
      mapped.owner = { type: 'personal' as DtoEditableOwnerModelTypeEnum };
    }

    // Map optional fields
    if (!this.isEmpty(rawRecord.friendlyName)) {
      mapped.friendlyName = String(rawRecord.friendlyName).trim();
    }

    if (!this.isEmpty(rawRecord.endDate)) {
      mapped.endDate = this.parseDate(String(rawRecord.endDate));
    }

    if (!this.isEmpty(rawRecord.customRecurrency)) {
      const recurrency = Number(rawRecord.customRecurrency);
      if (!isNaN(recurrency)) {
        mapped.customRecurrency = recurrency;
      }
    }

    // Map custom price
    if (rawRecord.customPrice && typeof rawRecord.customPrice === 'object') {
      mapped.customPrice = this.mapCustomPrice(rawRecord.customPrice);
    } else if (!this.isEmpty(rawRecord.customPriceAmount) && !this.isEmpty(rawRecord.customPriceCurrency)) {
      mapped.customPrice = this.mapCustomPrice({
        value: rawRecord.customPriceAmount,
        currency: rawRecord.customPriceCurrency,
      });
    }

    // Map labels (array of strings)
    if (rawRecord.labels) {
      if (Array.isArray(rawRecord.labels)) {
        mapped.labels = rawRecord.labels.map((label: any) => String(label).trim()).filter(Boolean);
      } else if (typeof rawRecord.labels === 'string') {
        // Support comma-separated string
        mapped.labels = rawRecord.labels.split(',').map((label: string) => label.trim()).filter(Boolean);
      }
    }

    // Map payer
    if (rawRecord.payer && typeof rawRecord.payer === 'object') {
      mapped.payer = this.mapPayer(rawRecord.payer);
    } else if (!this.isEmpty(rawRecord.payerType) && !this.isEmpty(rawRecord.payerFamilyId)) {
      mapped.payer = this.mapPayer({
        type: rawRecord.payerType,
        familyId: rawRecord.payerFamilyId,
        memberId: rawRecord.payerMemberId,
      });
    }

    // Map free trial
    if (rawRecord.freeTrial && typeof rawRecord.freeTrial === 'object') {
      mapped.freeTrial = this.mapFreeTrial(rawRecord.freeTrial);
    } else if (!this.isEmpty(rawRecord.freeTrialStartDate) && !this.isEmpty(rawRecord.freeTrialEndDate)) {
      mapped.freeTrial = this.mapFreeTrial({
        startDate: rawRecord.freeTrialStartDate,
        endDate: rawRecord.freeTrialEndDate,
      });
    }

    // Map family users (array of strings)
    if (rawRecord.familyUsers) {
      if (Array.isArray(rawRecord.familyUsers)) {
        mapped.familyUsers = rawRecord.familyUsers.map((user: any) => String(user).trim()).filter(Boolean);
      } else if (typeof rawRecord.familyUsers === 'string') {
        // Support comma-separated string
        mapped.familyUsers = rawRecord.familyUsers.split(',').map((user: string) => user.trim()).filter(Boolean);
      }
    }

    return mapped;
  }

  validate(record: Partial<DtoCreateSubscriptionRequest>): ValidationResult {
    const errors: ValidationError[] = [];

    // Validate required fields
    if (this.isEmpty(record.providerId)) {
      errors.push(this.createError('providerId', 'Provider ID is required'));
    }

    if (!record.startDate) {
      errors.push(this.createError('startDate', 'Start date is required'));
    } else if (!(record.startDate instanceof Date) || isNaN(record.startDate.getTime())) {
      errors.push(this.createError('startDate', 'Start date must be a valid date'));
    }

    if (this.isEmpty(record.recurrency)) {
      errors.push(this.createError('recurrency', 'Recurrency is required'));
    } else {
      const validRecurrencies = ['daily', 'weekly', 'monthly', 'quarterly', 'yearly', 'custom'];
      if (!validRecurrencies.includes(record.recurrency!.toLowerCase())) {
        errors.push(
          this.createError(
            'recurrency',
            `Recurrency must be one of: ${validRecurrencies.join(', ')}`
          )
        );
      }
    }

    if (!record.owner) {
      errors.push(this.createError('owner', 'Owner is required'));
    } else {
      const ownerErrors = this.validateOwner(record.owner);
      errors.push(...ownerErrors);
    }

    // Validate optional fields
    if (record.endDate) {
      if (!(record.endDate instanceof Date) || isNaN(record.endDate.getTime())) {
        errors.push(this.createError('endDate', 'End date must be a valid date'));
      } else if (record.startDate && record.endDate < record.startDate) {
        errors.push(this.createError('endDate', 'End date must be after start date'));
      }
    }

    if (record.customRecurrency !== undefined) {
      if (typeof record.customRecurrency !== 'number' || record.customRecurrency <= 0) {
        errors.push(this.createError('customRecurrency', 'Custom recurrency must be a positive number'));
      }
    }

    // Validate custom price
    if (record.customPrice) {
      const priceErrors = this.validateCustomPrice(record.customPrice);
      errors.push(...priceErrors);
    }

    // Validate payer
    if (record.payer) {
      const payerErrors = this.validatePayer(record.payer);
      errors.push(...payerErrors);
    }

    // Validate free trial
    if (record.freeTrial) {
      const freeTrialErrors = this.validateFreeTrial(record.freeTrial);
      errors.push(...freeTrialErrors);
    }

    return {
      isValid: errors.length === 0,
      errors,
    };
  }

  private mapOwner(rawOwner: Record<string, any>): DtoEditableOwnerModel {
    const owner: DtoEditableOwnerModel = {
      type: (rawOwner.type || 'personal') as DtoEditableOwnerModelTypeEnum,
    };

    if (rawOwner.familyId) {
      owner.familyId = String(rawOwner.familyId);
    }

    return owner;
  }

  private mapCustomPrice(rawPrice: Record<string, any>): DtoAmountModel | undefined {
    if (this.isEmpty(rawPrice.value) || this.isEmpty(rawPrice.currency)) {
      return undefined;
    }

    const value = Number(rawPrice.value);
    if (isNaN(value)) {
      return undefined;
    }

    return {
      value,
      currency: String(rawPrice.currency).toUpperCase(),
    };
  }

  private mapPayer(rawPayer: Record<string, any>): DtoEditableSubscriptionPayerModel | undefined {
    if (this.isEmpty(rawPayer.type) || this.isEmpty(rawPayer.familyId)) {
      return undefined;
    }

    const payer: DtoEditableSubscriptionPayerModel = {
      type: rawPayer.type as DtoEditableSubscriptionPayerModelTypeEnum,
      familyId: String(rawPayer.familyId),
    };

    if (rawPayer.memberId) {
      payer.memberId = String(rawPayer.memberId);
    }

    return payer;
  }

  private mapFreeTrial(rawFreeTrial: Record<string, any>): DtoSubscriptionFreeTrialModel | undefined {
    const startDate = this.parseDate(String(rawFreeTrial.startDate));
    const endDate = this.parseDate(String(rawFreeTrial.endDate));

    if (!startDate || !endDate) {
      return undefined;
    }

    return {
      startDate,
      endDate,
    };
  }

  private validateOwner(owner: DtoEditableOwnerModel): ValidationError[] {
    const errors: ValidationError[] = [];

    if (!owner.type) {
      errors.push(this.createError('owner.type', 'Owner type is required'));
    } else if (!['personal', 'family', 'system'].includes(owner.type)) {
      errors.push(this.createError('owner.type', 'Owner type must be one of: personal, family, system'));
    }

    if (owner.type === 'family' && this.isEmpty(owner.familyId)) {
      errors.push(this.createError('owner.familyId', 'Family ID is required when owner type is family'));
    }

    return errors;
  }

  private validateCustomPrice(price: DtoAmountModel): ValidationError[] {
    const errors: ValidationError[] = [];

    if (this.isEmpty(price.value)) {
      errors.push(this.createError('customPrice.value', 'Price value is required'));
    } else if (typeof price.value !== 'number' || price.value < 0) {
      errors.push(this.createError('customPrice.value', 'Price value must be a non-negative number'));
    }

    if (this.isEmpty(price.currency)) {
      errors.push(this.createError('customPrice.currency', 'Price currency is required'));
    } else if (price.currency.length !== 3) {
      errors.push(this.createError('customPrice.currency', 'Currency must be a 3-letter ISO code (e.g., USD, EUR)'));
    }

    return errors;
  }

  private validatePayer(payer: DtoEditableSubscriptionPayerModel): ValidationError[] {
    const errors: ValidationError[] = [];

    if (!payer.type) {
      errors.push(this.createError('payer.type', 'Payer type is required'));
    } else if (!['family', 'family_member'].includes(payer.type)) {
      errors.push(this.createError('payer.type', 'Payer type must be one of: family, family_member'));
    }

    if (this.isEmpty(payer.familyId)) {
      errors.push(this.createError('payer.familyId', 'Payer family ID is required'));
    }

    if (payer.type === 'family_member' && this.isEmpty(payer.memberId)) {
      errors.push(this.createError('payer.memberId', 'Member ID is required when payer type is family_member'));
    }

    return errors;
  }

  private validateFreeTrial(freeTrial: DtoSubscriptionFreeTrialModel): ValidationError[] {
    const errors: ValidationError[] = [];

    if (!freeTrial.startDate) {
      errors.push(this.createError('freeTrial.startDate', 'Free trial start date is required'));
    } else if (!(freeTrial.startDate instanceof Date) || isNaN(freeTrial.startDate.getTime())) {
      errors.push(this.createError('freeTrial.startDate', 'Free trial start date must be a valid date'));
    }

    if (!freeTrial.endDate) {
      errors.push(this.createError('freeTrial.endDate', 'Free trial end date is required'));
    } else if (!(freeTrial.endDate instanceof Date) || isNaN(freeTrial.endDate.getTime())) {
      errors.push(this.createError('freeTrial.endDate', 'Free trial end date must be a valid date'));
    }

    if (
      freeTrial.startDate &&
      freeTrial.endDate &&
      freeTrial.startDate instanceof Date &&
      freeTrial.endDate instanceof Date &&
      !isNaN(freeTrial.startDate.getTime()) &&
      !isNaN(freeTrial.endDate.getTime()) &&
      freeTrial.endDate < freeTrial.startDate
    ) {
      errors.push(this.createError('freeTrial.endDate', 'Free trial end date must be after start date'));
    }

    return errors;
  }
}
