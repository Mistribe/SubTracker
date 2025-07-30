// ignore_for_file: type=lint
/// auto generated
///  @Description Billing recurrency pattern (monthly, yearly, custom, etc.)
enum SubscriptionModelRecurrency {
    daily('daily'),
    weekly('weekly'),
    monthly('monthly'),
    quarterly('quarterly'),
    yearly('yearly'),
    custom('custom');
    const SubscriptionModelRecurrency(this.value);
    final String value;
}
