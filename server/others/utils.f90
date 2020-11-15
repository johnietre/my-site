module utils
contains

subroutine calc_sma(arr, prices, period, len)
  integer :: period, len, i
  real :: arr(len), prices(len), su = 0
  do i=1, len
    if (i < period) then
      su = su + prices(i)
    else
      su = su + prices(i) - prices(i-period)
    end if
    if (i < period) then
      arr(i) = 0
    else
      arr(i) = su / period
    endif
  end do
end subroutine calc_sma

subroutine calc_ema(arr, prices, period, len)
  integer :: period, len
  real :: arr(len), prices(len)
  call calc_sma(arr, prices, period, period)
  do i=period+1, len
    arr(i) = (prices(i) * (2.0 / (1.0 + period))) + (arr(i-1) * (1.0 - (2.0 / (1.0 + period))))
  end do
end subroutine calc_ema

subroutine calc_weighted(arr, prices, period, len)
  integer :: period, len, i, j
  real :: arr(len), prices(len), sum
  do i=1, len
    if (i < period) then
      arr(i) = 0
    else
      su = 0
      do j=period, 0, -1
        su = su + prices(i-j) * (period - j)
      end do
      arr(i) = su / (period * (period + 1.0) / 2.0)
      arr(i) = su / (period * (period + 1.0) / 2.0)
    end if
  end do
end subroutine calc_weighted

subroutine calc_wilders(arr, prices, period, len)
  integer :: period, len, i
  real :: arr(len), prices(len), sum
  call calc_sma(arr, prices, period, period)
  do i=period+1, len
    arr(i) = (prices(i) * (1.0 / period)) + (arr(i-1) * (1.0 - (1.0 / period)))
  end do
  return
end subroutine calc_wilders

subroutine calc_hull(arr, prices, period, len)
  integer :: period, len
  real :: arr(len), prices(len)
  prices(period) = prices(period)
  arr(len) = arr(len)
  return
end subroutine calc_hull

end module utils