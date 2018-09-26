package com.android.productchange.interfaces;

import com.android.productchange.adapters.ListViewProductReturnBooksAdapter;

import java.util.ArrayList;

/**
 * Created by cong-pv
 * on 2018/09/05.
 */

public interface AsyncResponseListReturnBooks {

    void processFinish(ArrayList output, ListViewProductReturnBooksAdapter apdaterReturnBooks);
}
