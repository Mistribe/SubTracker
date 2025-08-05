// ignore_for_file: type=lint
/// auto generated
///  @Description Billing recurrency pattern (monthly, yearly, custom, etc.)
enum SubscriptionModelRecurrency {
    unknown('unknown'),
    oneTime('one_time'),
    monthly('monthly'),
    quarterly('quarterly'),
    halfYearly('half_yearly'),
    yearly('yearly'),
    custom('custom');
    const SubscriptionModelRecurrency(this.value);
    final String value;
}
