import 'package:subscription_tracker/models/subscription.dart';

class Paginated<T> {
  final List<T> data;
  final int length;
  final int total;

  static Paginated<T> empty<T>() => Paginated<T>(null, null, null);

  Paginated(List<T>? data, int? length, int? total)
    : data = data ?? [],
      length = length ?? 0,
      total = total ?? 0;
}
